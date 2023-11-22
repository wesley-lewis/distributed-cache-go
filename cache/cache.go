package cache

import (
	"sync"
	"time"
	"fmt"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache {
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key[]byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock();
	defer c.lock.RUnlock()

	keyStr := string(key)
	_, ok := c.data[keyStr]
	return ok
}

func(c *Cache) Get(key []byte) ([]byte, error ) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	
	keyStr := string(key)
	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("Key %s doesn't exist", keyStr)
	}
	
	fmt.Printf("GET %s: %s", keyStr, string(val))
	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	
	go func() {
		<- time.After(ttl)
		delete(c.data, string(key))
	}()

	c.data[string(key)] = value
	fmt.Printf("SET %s to %s\n", string(key), string(value))
	return nil
}
