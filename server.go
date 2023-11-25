package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/wesley-lewis/distributed-cache/cache"
	"github.com/wesley-lewis/distributed-cache/proto"
)

type ServerOpts struct {
	ListenAddr string 
	IsLeader bool
	LeaderAddr string 
}

type Server struct {
	ServerOpts

	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher ) *Server {
	return &Server {
		ServerOpts: opts,
		cache: c,
		// TODO: only allocate this when we are the leader
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
	defer conn.Close()
	
	fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed:", conn.RemoteAddr())
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
}

func(s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
		case *proto.CommandSet:
		s.handleSetCommand(conn, v)
		
		case *proto.CommandGet:
			s.handleGetCommand(conn, v)
	}
}

func(s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)
	return s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL))
}

func(s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	return nil
}
