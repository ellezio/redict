package redict

import (
	"errors"
	"fmt"
)

type Database struct {
	store map[string]any
}

func NewDatabase() *Database {
	s := make(map[string]any)
	return &Database{
		store: s,
	}
}

func (db *Database) Set(key string, value []byte) error {
	s, ok := db.store[key]
	if !ok {
		strStore := newStrings()
		strStore.set(value)
		db.store[key] = strStore
		return nil
	}

	if strStore, ok := s.(*string_); ok {
		strStore.set(value)
		return nil
	} else {
		return fmt.Errorf("key %q already exists and is not type of String", key)
	}
}

func (db *Database) Get(key string) ([]byte, error) {
	s, ok := db.store[key]
	if !ok {
		return nil, fmt.Errorf("key %q does not exist", key)
	}

	if strStore, ok := s.(*string_); ok {
		return strStore.get(), nil
	} else {
		return nil, errors.New("invalid opertaion on String store")
	}
}
