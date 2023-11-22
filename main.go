package main

import (
	"net"
	"log"
	"time"

	"github.com/wesley-lewis/distributed-cache/cache"
)


func main() {
	opts := ServerOpts {
		ListenAddr: ":3000",
		IsLeader: true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp",":3000")
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Write([]byte("SET Foo Bar 34000000000"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d bytes\n", n)

		time.Sleep(time.Second * 2)
		conn.Write([]byte("GET Foo"))
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d bytes\n", n)
		log.Println("Value:", string(buf[:n]))
	}()

	server := NewServer(opts, cache.New())
	server.Start()
}
