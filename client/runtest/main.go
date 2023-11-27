package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wesley-lewis/distributed-cache/client"
)

func main() {
	sendStuff()
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
