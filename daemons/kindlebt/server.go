package main

import (
	"context"
	"log"

	kbt "github.com/Sighery/gokindlebt"
	pb "github.com/Sighery/turnkey/daemons/protogen"
)

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
