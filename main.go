package main

import (
	"context"
	"log"
	"fmt"
	"flag"
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
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			panic(err)
		}
		err = client.Set(context.Background(), []byte("foo"), []byte("bar"), 0)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 2)
		value, err := client.Get(context.Background(), []byte("foo"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Value: %s\n", value)
		
		client.Close()
	}()
	server := NewServer(opts, cache.New())
	server.Start()
}

