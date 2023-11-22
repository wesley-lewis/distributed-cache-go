package main

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"net"
	"strings"

	"github.com/wesley-lewis/distributed-cache/cache"
)

type ServerOpts struct {
	ListenAddr string 
	IsLeader bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher ) *Server {
	return &Server {
		ServerOpts: opts,
		cache: c,
	}
}

func(s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err.Error())
	}
	log.Printf("Server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err.Error())
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func(s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	rawStr := string(rawCmd)
	parts := strings.Split(rawStr, " ")
	if len(parts) == 0 {
		// respond
		log.Println("invalid command")
		conn.Write([]byte("invalid command"))
	}

	cmd := Command(parts[0])

	if cmd == CMDSet {
		if len(parts) != 4 {
			// respond
			log.Println("Invalid SET command")
			return 
		}
		// TODO: Need to check the error
		ttl, _ := strconv.Atoi(parts[3])
		key := []byte(parts[1])		
		value := []byte(parts[2])
		TTL := time.Duration(ttl)

		if err := s.handleSetCmd(conn, key, value, TTL); err != nil {
			// respond	
			return 
		}
	}
}

func (s *Server) handleSetCmd(conn net.Conn, key, value []byte, ttl time.Duration) error {
	
	return nil
}
