package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/wesley-lewis/distributed-cache/cache"
	"github.com/wesley-lewis/distributed-cache/client"
)


func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	opts := ServerOpts {
		ListenAddr: *listenAddr,
		IsLeader: len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2)
		for i:=0; i < 10; i++{
			SendCommand()
			time.Sleep(time.Millisecond * 200)
		}
		SendCommand()
	}()
	server := NewServer(opts, cache.New())
	server.Start()
}

// Just for testing purposes
func SendCommand() {
	cmd := &CommandSet{
		Key: []byte("foo"),
		Value: []byte("bar"),
		TTL: 0,
	}

	client, err := client.New(":3000", client.Options{})
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	n,_ := conn.Write(cmd.Bytes())
	fmt.Printf("Wrote %d bytes\n", n)
}
