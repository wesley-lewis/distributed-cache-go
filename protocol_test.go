package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	cmd := &CommandSet{
		Key: []byte("foo"),
		Value: []byte("bar"),
		TTL: 2,
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}

