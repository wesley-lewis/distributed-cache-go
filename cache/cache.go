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
	
	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	
	c.data[string(key)] = value

	if ttl > 0 {
		go func() {
			<- time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	return nil
}
