package flashdb

import (
	"reflect"
	"time"
)

func (db *FlashDB) Set(key string, value string) error {

	e := newRecord([]byte(key), []byte(value), StringRecord, StringSet)
	if err := db.write(e); err != nil {
		return err
	}

	err := db.set(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (db *FlashDB) SetEx(key string, value string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	db.strStore.Lock()
	defer db.strStore.Unlock()

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), []byte(value), ttl, StringRecord, StringExpire)
	if err := db.write(e); err != nil {
		return err
	}

	if err = db.set(key, value); err != nil {
		return
	}

	// set expired info.
	db.setTTL(String, key, ttl)
	return
}

func (db *FlashDB) Get(key string) (val string, err error) {
	db.strStore.RLock()
	defer db.strStore.RUnlock()

	val, err = db.get(key)
	if err != nil {
		return
	}

	return
}

func (db *FlashDB) Delete(key string) error {
	db.strStore.Lock()
	defer db.strStore.Unlock()

	e := newRecord([]byte(key), nil, StringRecord, StringRem)
	if err := db.write(e); err != nil {
		return err
	}

	db.strStore.Delete(key)
	db.exps.HDel(String, key)
	return nil
}

func (db *FlashDB) Expire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	db.strStore.Lock()
	defer db.strStore.Unlock()

	if _, err = db.get(key); err != nil {
		return
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, StringRecord, StringExpire)
	if err := db.write(e); err != nil {
		return err
	}

	db.setTTL(String, key, ttl)
	return
}

func (db *FlashDB) TTL(key string) (ttl int64) {

	db.strStore.Lock()
	defer db.strStore.Unlock()

	deadline := db.getTTL(String, key)
	if deadline == nil {
		return
	}

	if db.hasExpired(key, String) {
		db.evict(key, String)
		return
	}

	return deadline.(int64) - time.Now().Unix()
}

func (db *FlashDB) Exists(key string) bool {
	db.strStore.RLock()
	defer db.strStore.RUnlock()

	_, err := db.strStore.get(key)
	if err != nil {
		if err == ErrExpiredKey {
			db.evict(key, String)
		}
		return false
	}

	return true
}

func (db *FlashDB) set(key string, value string) error {
	var existVal string
	existVal, err := db.get(key)
	if err != nil && err != ErrExpiredKey && err != ErrInvalidKey {
		return err
	}

	if reflect.DeepEqual(existVal, value) {
		return err
	}

	db.strStore.Set(key, value)

	return nil
}

func (db *FlashDB) get(key string) (val string, err error) {
	v, err := db.strStore.get(key)
	if err != nil {
		return "", err
	}

	// Check if the key is expired.
	if db.hasExpired(key, String) {
		db.evict(key, String)
		return "", ErrExpiredKey
	}

	val = v.(string)
	return
}
