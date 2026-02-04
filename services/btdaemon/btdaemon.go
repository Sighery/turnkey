package btdaemon

import (
	"context"
	"log"

	pb "github.com/Sighery/turnkey/daemons/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotifySubscription struct {
	C    <-chan []byte
	Stop func()
}

type BtdaemonClient struct {
	conn    *grpc.ClientConn
	service pb.DaemonClient
}

func NewBtdaemonClient(addr string) (*BtdaemonClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &BtdaemonClient{}, err
	}

	daemon := pb.NewDaemonClient(conn)

	return &BtdaemonClient{conn: conn, service: daemon}, nil
}

func (c *BtdaemonClient) Close() error {
	return c.conn.Close()
}

func (c *BtdaemonClient) IsReady(ctx context.Context) error {
	r, err := c.service.IsReady(ctx, &pb.IsReadyRequest{})
	if err != nil {
		return err
	}

	log.Printf("BT daemon status: %t\n", r.GetStatus())
	return nil
}

func (c *BtdaemonClient) Connect(ctx context.Context, addr string) (string, error) {
	r, err := c.service.ConnectBle(ctx, &pb.ConnectBleRequest{Addr: addr})
	if err != nil {
		return "", err
	}

	return r.GetConnectionId(), nil
}

func (c *BtdaemonClient) ReadChar(ctx context.Context, connId string, charId string) (
	[]byte, error,
) {
	r, err := c.service.ReadChar(
		ctx, &pb.ReadCharRequest{ConnectionId: connId, CharacteristicId: charId},
	)
	if err != nil {
		return nil, err
	}

	return r.GetResponse(), nil
}

func (c *BtdaemonClient) WriteChar(
	ctx context.Context, connId string, charId string, withResponse bool, data []byte,
) ([]byte, error) {
	body := &pb.WriteCharRequest{
		ConnectionId:     connId,
		CharacteristicId: charId,
		WithResponse:     withResponse,
		Data:             data,
	}

	r, err := c.service.WriteChar(ctx, body)
	if err != nil {
		return nil, err
	}

	return r.GetResponse(), nil
}

func (c *BtdaemonClient) NotifyChar(ctx context.Context, connId string, charId string) (
	*NotifySubscription, error,
) {
	cancelCtx, cancel := context.WithCancel(ctx)

	body := &pb.NotifyCharRequest{ConnectionId: connId, CharacteristicId: charId}
	stream, err := c.service.NotifyChar(cancelCtx, body)
	if err != nil {
		cancel()
		return nil, err
	}

	out := make(chan []byte)

	go func() {
		defer close(out)
		for {
			resp, err := stream.Recv()
			if err != nil {
				cancel()
				return
			}

			out <- resp.GetResponse()
		}
	}()

	return &NotifySubscription{
		C: out,
		Stop: func() {
			cancel()
		},
	}, nil
}
