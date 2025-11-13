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

func (db *Database) LPush(key string, value []byte) error {
	s, ok := db.store[key]
	if !ok {
		l := newList()
		l.insertHead(value)
		db.store[key] = l
		return nil
	}

	if l, ok := s.(*list); ok {
		l.insertHead(value)
		return nil
	} else {
		return errors.New("wrong operation against key")
	}
}

func (db *Database) RPush(key string, value []byte) error {
	s, ok := db.store[key]
	if !ok {
		l := newList()
		l.insertTail(value)
		db.store[key] = l
		return nil
	}

	if l, ok := s.(*list); ok {
		l.insertTail(value)
		return nil
	} else {
		return errors.New("wrong operation against key")
	}
}

func (db *Database) LPop(key string) ([]byte, error) {
	s, ok := db.store[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if l, ok := s.(*list); ok {
		b := l.popHead()
		return b, nil
	} else {
		return nil, errors.New("wrong operation against key")
	}
}

func (db *Database) RPop(key string) ([]byte, error) {
	s, ok := db.store[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if l, ok := s.(*list); ok {
		b := l.popTail()
		return b, nil
	} else {
		return nil, errors.New("wrong operation against key")
	}
}

func (db *Database) LRange(key string, start, end int64) ([][]byte, error) {
	s, ok := db.store[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if l, ok := s.(*list); ok {
		return l.get(start, end), nil
	} else {
		return nil, errors.New("wrong operation against key")
	}
}

func (db *Database) LTrim(key string, start, end int64) error {
	s, ok := db.store[key]
	if !ok {
		return errors.New("key does not exist")
	}

	if l, ok := s.(*list); ok {
		l.trim(start, end)
		return nil
	} else {
		return errors.New("wrong operation against key")
	}
}

func (db *Database) LLen(key string) (uint32, error) {
	s, ok := db.store[key]
	if !ok {
		return 0, errors.New("key does not exist")
	}

	if l, ok := s.(*list); ok {
		return l.length, nil
	} else {
		return 0, errors.New("wrong operation against key")
	}
}
