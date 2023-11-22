package main

import (
	"net"
	"time"
	"log"

	"github.com/wesley-lewis/distributed-cache/cache"
)


func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}

	go func() {
		for {
			time.Sleep(time.Second * 2)
			conn, err := net.Dial("tcp",":3000")
			if err != nil {
				log.Fatal(err)
			}

			n, err := conn.Write([]byte("SET Foo Bar 3400"))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Wrote %d bytes\n", n)
		}
	}()
	
	server := NewServer(opts, cache.New())
	server.Start()
}
