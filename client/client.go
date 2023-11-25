package client

import (
	"context"
	"net"

	"github.com/wesley-lewis/distributed-cache/proto"
)

type Options struct {

}

type Client struct{
	conn net.Conn
	opts Options
}

func New(endpoint string , opts Options) (*Client, error){
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client {
		conn: conn,
		opts: opts,
	},nil
}

func(c *Client) Set(ctx context.Context, key, value []byte) (any, error) {
	cmd := &proto.CommandSet {
		Key: key,
		Value: value,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func(c *Client) Close() error {
	return c.conn.Close()
}
