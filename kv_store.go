package main

import "errors"

var store = make(map[string]string)

var ErrorKeyNotFound = errors.New("Key not found")

// will overwrite existing value
func Put(key, value string) error {
	store[key] = value
	return nil
}

func Get(key string) (string, error) {
	value, exists := store[key]

	if !exists {
		return "", ErrorKeyNotFound // "" is the zero value for string
	}

	return value, nil
}

func Delete(key string) error {
	delete(store, key)
	return nil
}