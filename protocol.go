package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Command byte 

const (
	CmdNonce Command = iota
	CmdSet   
	CmdGet
	CmdDel
)

type CommandSet struct {
	Key []byte 
	Value []byte 
	TTL int  
}

func(c *CommandSet) Bytes() []byte{
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)
	binary.Write(buf, binary.LittleEndian, int64(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)	

	binary.Write(buf, binary.LittleEndian, int64(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)	

	binary.Write(buf, binary.LittleEndian, c.TTL)	

	return buf.Bytes()

}

func ParseCommand(r io.Reader) {

}
