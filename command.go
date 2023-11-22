package main 

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"log"
)

type Command string 

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type MSGSet struct {
	Key []byte 
	Value []byte 
	TTL time.Duration
}

type MSGGet struct {
	Key []byte 
}

type Message struct {
	Cmd Command 
	Key []byte 
	Value []byte 
	TTL time.Duration
}

func(m *Message) ToBytes() []byte{
	switch m.Cmd {
	case CMDGet:
		cmd := fmt.Sprintf("%s %s",m.Cmd, m.Key)
		return []byte(cmd)
	case CMDSet:
		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL)
		return []byte(cmd)
	default: 
		panic("unknown command")
	}
}

func parseCommand(raw []byte) (*Message, error) {
	var (
		rawStr = string(raw)
		parts = strings.Split(rawStr, " ")
	)
	if len(parts) < 2 {
		// respond
		log.Println("invalid command")
		return nil, fmt.Errorf("invalid protocol format")
	}

	msg := &Message {
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}
	
	if msg.Cmd == CMDSet {
		if len(parts) != 4 {
			return nil, errors.New("invalid SET command")
		}
		msg.Value = []byte(parts[2])
		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid TTL field")
		}
		msg.TTL = time.Duration(ttl)
		
	}

	return msg, nil
}
