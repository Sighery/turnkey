package main

import (
	"flag"
	"log"
	"net"

	kbt "github.com/Sighery/gokindlebt"
	pb "github.com/Sighery/turnkey/daemons/protogen"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "0.0.0.0:50010", "ip:port to use")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	if err := kbt.UseBluetoothPrivileges(); err != nil {
		log.Fatalf("Failed to drop Bluetooth privileges %v", err)
	}

	service, err := NewDaemonService()
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterDaemonServer(server, &service)

	log.Printf("server listening at %v", listener.Addr())
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
