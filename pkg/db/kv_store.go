package db

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	data map[string]string
}{data: make(map[string]string)}

var ErrorKeyNotFound = errors.New("key not found")

func Put(key, value string) error {
	store.Lock()
	defer store.Unlock()

	store.data[key] = value
	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	value, exists := store.data[key]

	if !exists {
		return "", ErrorKeyNotFound
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.data, key)
	return nil
}

func GetAll() map[string]string {
	return store.data
}
