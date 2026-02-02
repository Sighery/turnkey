package btdaemon

import (
	"context"
	"log"

	pb "github.com/Sighery/turnkey/daemons/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	log.Printf("BT daemon status: %t", r.GetStatus())
	return nil
}
