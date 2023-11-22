package main 

import (
	"time"
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

