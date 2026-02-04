package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"

	kbt "github.com/Sighery/gokindlebt"
	pb "github.com/Sighery/turnkey/daemons/protogen"
)

var (
	globalCounter atomic.Uint64
	connections   sync.Map
)

type Connection struct {
	id      string
	adapter kbt.Adapter
	session kbt.Session
	conn    kbt.BleConnection
}

func NewConnection(adapter kbt.Adapter, session kbt.Session, conn kbt.BleConnection) Connection {
	currentCount := globalCounter.Add(1)
	id := strconv.FormatUint(currentCount, 10)
	connection := Connection{id: id, adapter: adapter, session: session, conn: conn}
	connections.Store(id, connection)
	return connection
}

func (c Connection) Close() error {
	errors := false

	err := c.adapter.DisconnectBle(c.conn)
	if err != nil {
		errors = true
		log.Println("ID:%s Failed to disconnect BLE device: %w", c.id, err)
	}

	err = c.adapter.DeregisterGattClient(c.session)
	if err != nil {
		errors = true
		log.Println("ID:%s Failed to deregister GATT Client: %w", c.id, err)
	}

	err = c.adapter.DeregisterBle(c.session)
	if err != nil {
		errors = true
		log.Println("ID:%s Failed to deregister BLE: %w", c.id, err)
	}

	err = c.adapter.CloseSession(c.session)
	if err != nil {
		errors = true
		log.Println("ID:%s Failed to close session: %w", c.id, err)
	}

	connections.Delete(c.id)

	if errors {
		return fmt.Errorf("Errors while trying to close KindleBT connection %s", c.id)
	}

	return nil
}

func getConnection(id string) (Connection, error) {
	v, exists := connections.Load(id)
	if !exists {
		return Connection{}, fmt.Errorf("No connection %s", id)
	}

	conn, ok := v.(Connection)
	if !ok {
		return Connection{}, fmt.Errorf("Couldn't cast connection %s: %v", id, conn)
	}

	return conn, nil
}

type DaemonService struct {
	adapter kbt.Adapter
	pb.UnimplementedDaemonServer
}

func NewDaemonService() (DaemonService, error) {
	adapter, err := kbt.NewAdapter()
	if err != nil {
		return DaemonService{}, err
	}

	return DaemonService{adapter: adapter}, nil
}

func (s *DaemonService) IsReady(_ context.Context, req *pb.IsReadyRequest) (
	*pb.IsReadyResponse, error,
) {
	log.Println("Received IsReady request")

	log.Println("Is ble supported check")
	if isBle := s.adapter.IsBleSupported(); isBle == false {
		log.Printf("BLE is not supported\n")
		return &pb.IsReadyResponse{Status: false}, nil
	}

	log.Println("Supported session check")
	if support := s.adapter.GetSupportedSession(); support == int(kbt.SessionNone) {
		log.Printf("Session is down %d\n", support)
		return &pb.IsReadyResponse{Status: false}, nil
	}

	return &pb.IsReadyResponse{Status: true}, nil
}

func (s *DaemonService) ConnectBle(_ context.Context, req *pb.ConnectBleRequest) (
	*pb.ConnectBleResponse, error,
) {
	log.Println("Received ConnectBle request")

	session, err := s.adapter.OpenSession(kbt.SessionDual)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't open session: %w", err)
	}

	err = s.adapter.RegisterBle(session)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't register BLE: %w", err)
	}

	err = s.adapter.RegisterGattClient(session)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't register GATT Client: %w", err)
	}

	addr, err := kbt.NewAddressFromString(req.GetAddr())
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't parse BLE address: %w", err)
	}

	conn, err := s.adapter.ConnectBleSimple(session, addr)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't connect to BLE device: %w", err)
	}

	err = s.adapter.DiscoverServices(session, conn)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't discover services: %w", err)
	}

	_, err = s.adapter.RetrieveGattDatabase(conn)
	if err != nil {
		return &pb.ConnectBleResponse{}, fmt.Errorf("Couldn't retrieve GATT DB: %w", err)
	}

	connection := NewConnection(s.adapter, session, conn)
	fmt.Println("Connid was", connection.id)
	return &pb.ConnectBleResponse{ConnectionId: connection.id}, nil
}

func (s *DaemonService) ReadChar(_ context.Context, req *pb.ReadCharRequest) (
	*pb.ReadCharResponse, error,
) {
	log.Println("Received ReadChar request")

	connId := req.GetConnectionId()
	charId := req.GetCharacteristicId()

	conn, err := getConnection(connId)
	if err != nil {
		return &pb.ReadCharResponse{}, fmt.Errorf("Connection %s not registered: %w", connId, err)
	}

	uuid, err := kbt.NewCharacteristicUuidFromString(charId)
	if err != nil {
		return &pb.ReadCharResponse{}, fmt.Errorf("Characteristic UUID invalid: %w", err)
	}

	value, err := s.adapter.ReadCharacteristic(conn.session, conn.conn, uuid)
	if err != nil {
		return &pb.ReadCharResponse{}, fmt.Errorf("Failed to read characteristic: %w", err)
	}

	blob, ok := value.Value.(kbt.CharacteristicValueBlob)
	if !ok {
		return &pb.ReadCharResponse{}, fmt.Errorf("Non-blob values not implemented")
	}

	return &pb.ReadCharResponse{Response: blob.V}, nil
}

func (s *DaemonService) WriteChar(_ context.Context, req *pb.WriteCharRequest) (
	*pb.WriteCharResponse, error,
) {
	log.Println("Received WriteChar request")

	connId := req.GetConnectionId()
	charId := req.GetCharacteristicId()
	withResponse := req.GetWithResponse()
	data := req.GetData()

	conn, err := getConnection(connId)
	if err != nil {
		return &pb.WriteCharResponse{}, fmt.Errorf("Connection %s not registered: %w", connId, err)
	}

	uuid, err := kbt.NewCharacteristicUuidFromString(charId)
	if err != nil {
		return &pb.WriteCharResponse{}, fmt.Errorf("Characteristic UUID invalid: %w", err)
	}

	value := kbt.CharacteristicValueBlob{V: data}

	if !withResponse {
		err = s.adapter.WriteCharacteristicWithoutResponse(conn.session, conn.conn, uuid, value)
		if err != nil {
			return &pb.WriteCharResponse{}, fmt.Errorf("Couldn't write characteristic: %w", err)
		}

		return &pb.WriteCharResponse{Response: []byte{}}, nil
	} else {
		res, err := s.adapter.WriteCharacteristicWithResponse(conn.session, conn.conn, uuid, value)
		if err != nil {
			return &pb.WriteCharResponse{}, fmt.Errorf("Couldn't write characteristic: %w", err)
		}

		blob, ok := res.Value.(kbt.CharacteristicValueBlob)
		if !ok {
			return &pb.WriteCharResponse{}, fmt.Errorf("Non-blob values not implemented")
		}

		return &pb.WriteCharResponse{Response: blob.V}, nil
	}
}
