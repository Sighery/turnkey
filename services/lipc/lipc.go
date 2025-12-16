package lipc

import (
	"context"

	"github.com/godbus/dbus/v5"

	lipcapi "github.com/clintharrison/liblipc-go/lipc"
)

type LipcClient struct {
	conn *dbus.Conn
}

func NewLipcClient(dbusConn *dbus.Conn) *LipcClient {
	return &LipcClient{conn: dbusConn}
}

func (c *LipcClient) GetString(
	ctx context.Context, service, property string,
) (string, error) {
	return lipcapi.LipcGetProperty[string](ctx, c.conn, service, property)
}

func (c *LipcClient) SetString(
	ctx context.Context, service, property string, value string,
) error {
	return lipcapi.LipcSetProperty[string](ctx, c.conn, service, property, value)
}

func (c *LipcClient) GetInt(
	ctx context.Context, service, property string,
) (int32, error) {
	return lipcapi.LipcGetProperty[int32](ctx, c.conn, service, property)
}

func (c *LipcClient) SetInt(
	ctx context.Context, service, property string, value int32,
) error {
	return lipcapi.LipcSetProperty[int32](ctx, c.conn, service, property, value)
}
