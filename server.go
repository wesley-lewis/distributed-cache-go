package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/wesley-lewis/distributed-cache/cache"
	"github.com/wesley-lewis/distributed-cache/client"
	"github.com/wesley-lewis/distributed-cache/proto"
	"go.uber.org/zap"
)

type ServerOpts struct {
	ListenAddr string 
	IsLeader bool
	LeaderAddr string 
}

type Server struct {
	ServerOpts

	members map[*client.Client]struct{}

	cache cache.Cacher
	logger *zap.SugaredLogger
}

func NewServer(opts ServerOpts, c cache.Cacher ) *Server {
	l,_ := zap.NewProduction()

	return &Server {
		ServerOpts: opts,
		cache: c,
		members: make(map[*client.Client]struct{}),
		logger: l.Sugar(),
	}
}

func(s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err.Error())
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			err := s.dialLeader()
			if err != nil {
				log.Println(err)
			}
		}()
		
	}
	
	s.logger.Infow("Server starting", "addr", s.ListenAddr, "leader", s.LeaderAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func(s *Server) dialLeader() error{
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("Failed to dial leader [%s]", s.LeaderAddr)
	}
	s.logger.Infow("connected to leader", "addr",s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)

	s.handleConn(conn)
	return nil
}

/// handleConn first parses the command i.e. the data passed in through the connection 
/// then forwards the command to the handler based on our protocol.
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	
	fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				// fmt.Println("connection closed:", conn.RemoteAddr())
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
}

/// handleCommadn acts as a router
/// It passes the connection and the command to the handler 
/// The connection can be used to send a response
func(s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
		case *proto.CommandSet:
		s.handleSetCommand(conn, v)
		
		case *proto.CommandGet:
			s.handleGetCommand(conn, v)
		case *proto.CommandJoin:
			s.handleJoinCommand(conn)
	}
}

func(s *Server) handleJoinCommand(conn net.Conn) error{
	fmt.Println("Member just joined", conn.RemoteAddr())

	c := client.NewFromConn(conn)
	s.members[c] = struct{}{}

	return nil
}

func(s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s\n", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := &proto.ResponseSet{
			// Status: proto.StatusError,
	}
	err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL))
	if err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = proto.StatusOK
	_, err = conn.Write(resp.Bytes())

	return err
}


func(s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet)  error {
	log.Printf("GET %s\n", cmd.Key)

	resp := proto.ResponseGet{}

	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusKeyNotFound
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = proto.StatusOK
	resp.Value = value 
	_, err = conn.Write(resp.Bytes())
	
	return err
}
