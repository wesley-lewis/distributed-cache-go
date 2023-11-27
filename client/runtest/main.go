package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/wesley-lewis/distributed-cache/client"
)

type Server struct {
	raft *raft.Raft
}

func main() {
	var (
		cfg				= raft.DefaultConfig()
		fsm				= &raft.MockFSM{}
		logStore		= &raft.InmemStore{}
		snapShotStore	= raft.NewInmemSnapshotStore()
		stableStore		= raft.NewInmemStore()
		timeout			= time.Duration(time.Second * 5)
	)
	
	cfg.LocalID = "ID"
	ips, err := net.LookupIP("localhost")
	if err != nil {
		log.Fatal(err)
	}
	if len(ips) == 0 {
		log.Fatalf("localhost did not resolve to any IPs")
	}
	addr := &net.TCPAddr{IP: ips[0], Port: 4000}

	tr, err := raft.NewTCPTransport(":4000", addr, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	r, err := raft.NewRaft(cfg, fsm, logStore, stableStore, snapShotStore, tr)	
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", r)
	select {
	}
}

func sendStuff() {
	c, err := client.New(":3000", client.Options{})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		var (
			key = []byte(fmt.Sprintf("key_%d", i))
			val = []byte(fmt.Sprintf("value_%d", i))
		)

		err = c.Set(context.Background(), key, val, 0)
		if err != nil {
			log.Fatal(err)
		}
		
		time.Sleep(time.Second)
	}
	c.Close()
}
