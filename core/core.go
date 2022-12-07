package core

import (
	"errors"
	"fmt"
	"sync"
)

type KVStore struct {
	sync.RWMutex
	data   map[string]string
	logger TransactionLogger // is "driven by" KVStore
}

func NewKVStore(logger TransactionLogger) *KVStore {
	return &KVStore{
		data:   make(map[string]string),
		logger: logger,
	}
}

var ErrorKeyNotFound = errors.New("key not found")

func (store *KVStore) Put(key, value string, isReplay bool) error {
	store.Lock()
	defer store.Unlock()

	store.data[key] = value
	if !isReplay {
		store.logger.LogPut(key, value)
	}
	return nil
}

func (store *KVStore) Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	value, exists := store.data[key]

	if !exists {
		return "", ErrorKeyNotFound
	}

	return value, nil
}

func (store *KVStore) Delete(key string, isReplay bool) error {
	store.Lock()
	defer store.Unlock()

	delete(store.data, key)
	if !isReplay {
		store.logger.LogDelete(key)
	}
	return nil
}

func (store *KVStore) GetAll() map[string]string {
	store.RLock()
	defer store.RUnlock()
	return store.data
}

func (store *KVStore) Restore() error {
	var err error

	fmt.Println("--> Replaying Logs...")
	replayEvents, replayErrors := store.logger.ReplayEvents()

	event, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-replayErrors: // got an error
			return fmt.Errorf("error replaying logs: %w", err)

		case event, ok = <-replayEvents: // got an event

			switch event.Type {
			case EventPut:
				err = store.Put(event.Key, event.Value, true)

			case EventDelete:
				err = store.Delete(event.Key, true)
			}
		}
	}

	fmt.Printf("--> Recreated DataStore: %v\n", store.GetAll())
	fmt.Println("--> Done replaying Logs")

	fmt.Println("--> Starting the Logger...")
	store.logger.Run()
	fmt.Println("--> Done starting the Logger")

	go func() {
		for err = range store.logger.Err() {
			fmt.Printf("logger error: %v", err)
		}
	}()

	return err
}

type EventType byte

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)

type Event struct {
	Index uint64    // index for ordering
	Type  EventType // type of event (Put, Delete, etc)
	Key   string    // key where event happened
	Value string    // value associated (with Put)
}

type TransactionLogger interface {
	LogPut(key, value string)
	LogDelete(key string)
	Err() <-chan error
	ReplayEvents() (<-chan Event, <-chan error)
	Run()
}
