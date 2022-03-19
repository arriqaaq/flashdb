package flashdb

import (
	"errors"
	"sync"
	"time"

	"github.com/arriqaaq/aol"
	"github.com/arriqaaq/hash"
)

var (
	ErrInvalidKey     = errors.New("invalid key")
	ErrInvalidTTL     = errors.New("invalid ttl")
	ErrExpiredKey     = errors.New("key has expired")
	ErrTxClosed       = errors.New("tx closed")
	ErrDatabaseClosed = errors.New("database closed")
	ErrTxNotWritable  = errors.New("tx not writable")
)

type (
	FlashDB struct {
		mu     sync.RWMutex
		config *Config
		exps   *hash.Hash // hashmap of ttl keys
		log    *aol.Log

		closed  bool // set when the database has been closed
		persist bool // do we write to disk

		strStore  *strStore
		hashStore *hashStore
		setStore  *setStore
		zsetStore *zsetStore

		evictors []evictor // background manager to delete keys periodically
	}
)

func New(config *Config) (*FlashDB, error) {

	config.validate()

	db := &FlashDB{
		config:    config,
		strStore:  newStrStore(),
		setStore:  newSetStore(),
		hashStore: newHashStore(),
		zsetStore: newZSetStore(),
		exps:      hash.New(),
	}

	evictionInterval := config.evictionInterval()
	if evictionInterval > 0 {
		db.evictors = []evictor{
			newSweeperWithStore(db.strStore, evictionInterval),
			newSweeperWithStore(db.setStore, evictionInterval),
			newSweeperWithStore(db.hashStore, evictionInterval),
			newSweeperWithStore(db.zsetStore, evictionInterval),
		}
		for _, evictor := range db.evictors {
			go evictor.run(db.exps)
		}
	}

	db.persist = config.Path != ""
	if db.persist {
		opts := aol.DefaultOptions
		opts.NoSync = config.NoSync

		l, err := aol.Open(config.Path, opts)
		if err != nil {
			return nil, err
		}

		db.log = l

		// load data from append-only log
		err = db.load()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *FlashDB) setTTL(dType DataType, key string, ttl int64) {
	db.exps.HSet(dType, key, ttl)
}

func (db *FlashDB) getTTL(dType DataType, key string) interface{} {
	return db.exps.HGet(dType, key)
}

func (db *FlashDB) hasExpired(key string, dType DataType) (expired bool) {
	ttl := db.exps.HGet(dType, key)
	if ttl == nil {
		return
	}
	if time.Now().Unix() > ttl.(int64) {
		expired = true
	}
	return
}

func (db *FlashDB) evict(key string, dType DataType) {
	ttl := db.exps.HGet(dType, key)
	if ttl == nil {
		return
	}

	var r *record
	if time.Now().Unix() > ttl.(int64) {
		switch dType {
		case String:
			r = newRecord([]byte(key), nil, StringRecord, StringRem)
			db.strStore.Delete([]byte(key))
		case Hash:
			r = newRecord([]byte(key), nil, HashRecord, HashHClear)
			db.hashStore.HClear(key)
		case Set:
			r = newRecord([]byte(key), nil, SetRecord, SetSClear)
			db.setStore.SClear(key)
		case ZSet:
			r = newRecord([]byte(key), nil, ZSetRecord, ZSetZClear)
			db.zsetStore.ZClear(key)
		}

		if err := db.write(r); err != nil {
			panic(err)
		}

		db.exps.HDel(dType, key)
	}
}

func (db *FlashDB) Close() error {
	db.closed = true
	for _, evictor := range db.evictors {
		evictor.stop()
	}
	if db.log != nil {
		err := db.log.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *FlashDB) write(r *record) error {
	if db.log == nil {
		return nil
	}
	encVal, err := r.encode()
	if err != nil {
		return err
	}

	return db.log.Write(encVal)
}
