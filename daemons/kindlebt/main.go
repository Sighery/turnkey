package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	kbt "github.com/Sighery/gokindlebt"
	pb "github.com/Sighery/turnkey/daemons/protogen"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "0.0.0.0:50010", "ip:port to use")
)

func freeConnections() {
	connections.Range(func(key, value any) bool {
		v, ok := value.(Connection)
		if !ok {
			fmt.Printf("Couldn't cast conn %s to type", key)
			return true
		}

		err := v.Close()
		if err != nil {
			fmt.Printf("Connection %s close wasn't clean", key)
		}

		return true
	})
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if err := kbt.UseBluetoothPrivileges(); err != nil {
		log.Fatalf("Failed to drop Bluetooth privileges %v", err)
	}

	service, err := NewDaemonService()
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterDaemonServer(server, &service)

	go func() {
		log.Printf("server listening at %v", listener.Addr())
		if err = server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	_ = <-sigs
	log.Println("Received shutdown signal. Shutting down gracefully...")
	freeConnections()

	server.GracefulStop()
	log.Println("Server shut down")
}
