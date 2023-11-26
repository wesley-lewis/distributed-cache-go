package client

import (
	"context"
	"fmt"
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

func(c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet {
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := proto.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status != proto.StatusOK {
		return nil, fmt.Errorf("[%s]\n", resp.Status.String())
	}

	return  resp.Value, nil
}

func(c *Client) Set(ctx context.Context, key, value []byte, ttl int) error {
	cmd := &proto.CommandSet {
		Key: key,
		Value: value,
		TTL: ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != proto.StatusOK {
		return fmt.Errorf("server responed with non OK status [%s]", string(resp.Status))
	}

	return  nil
}

func(c *Client) Close() error {
	return c.conn.Close()
}
