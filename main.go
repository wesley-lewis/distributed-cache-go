package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
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

	go multipleClients()
	server := NewServer(opts, cache.New())
	server.Start()
}

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	io.ReadFull(rand.Reader, buf)
	return buf
}

func multipleClients() {
	time.Sleep(time.Second * 2)
	for i := 0; i < 1000; i++ {
		go func() {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				panic(err)
			}

			var (
				key = randomBytes(10)
				val = randomBytes(10)
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
		}()
	}
}
