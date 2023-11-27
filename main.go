package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
		time.Sleep(time.Second * 10)
		if opts.IsLeader {
			multipleClients()
		}
	}()
	// go multipleClients()
	server := NewServer(opts, cache.New())
	server.Start()
}

func multipleClients() {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				panic(err)
			}

			var (
				key = []byte(fmt.Sprintf("key_%d", i))
				val = []byte(fmt.Sprintf("value_%d",i))
			)

			resp, err := client.Get(context.Background(), key)
			if err != nil {
				log.Printf("Error: %s", err)
			}

			fmt.Println(resp)

			err = client.Set(context.Background(), key, val, 0)
			if err != nil {
				log.Fatal(err)
			}

			value, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Value: %s\n", value)

			client.Close()
		}(i)
	}
}
// 1:13:47
