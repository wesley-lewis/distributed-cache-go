package main

import (
	"context"
	"fmt"
	"log"
	"net"

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
	msg, err := parseCommand(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command", err)
		// respond
		return 
	}
	
	switch msg.Cmd {
	case CMDSet: 
		 err = s.handleSetCmd(conn, msg) 
	
	case CMDGet: 
		var value []byte
		value, err = s.handleGetCmd(conn, msg)
		conn.Write(value)
	}

	if err != nil {
		fmt.Println("Failed to handle command:", err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server)handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO())
	return nil
}

func(s *Server) handleGetCmd(conn net.Conn, msg *Message) ([]byte, error) {
	value, err := s.cache.Get(msg.Key) 
	if err != nil {
		return nil, err
	}
	return value, nil
}

func(s *Server) sendToFollowers(ctx context.Context) error {
	return nil	
}
