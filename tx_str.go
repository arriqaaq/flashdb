package flashdb

import (
	"reflect"
	"time"
)

func (tx *Tx) Set(key string, value string) error {

	e := newRecord([]byte(key), []byte(value), StringRecord, StringSet)
	tx.addRecord(e)

	err := tx.set(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (tx *Tx) SetEx(key string, value string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), []byte(value), ttl, StringRecord, StringExpire)
	tx.addRecord(e)

	if err = tx.set(key, value); err != nil {
		return
	}

	// set expired info.
	tx.db.setTTL(String, key, ttl)
	return
}

func (tx *Tx) Get(key string) (val string, err error) {
	val, err = tx.get(key)
	if err != nil {
		return
	}

	return
}

func (tx *Tx) Delete(key string) error {

	e := newRecord([]byte(key), nil, StringRecord, StringRem)
	tx.addRecord(e)

	tx.db.strStore.Delete(key)
	tx.db.exps.HDel(String, key)
	return nil
}

func (tx *Tx) Expire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	if _, err = tx.get(key); err != nil {
		return
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, StringRecord, StringExpire)
	tx.addRecord(e)

	tx.db.setTTL(String, key, ttl)
	return
}

func (tx *Tx) TTL(key string) (ttl int64) {

	deadline := tx.db.getTTL(String, key)
	if deadline == nil {
		return
	}

	if tx.db.hasExpired(key, String) {
		tx.db.evict(key, String)
		return
	}

	return deadline.(int64) - time.Now().Unix()
}

func (tx *Tx) Exists(key string) bool {
	_, err := tx.db.strStore.get(key)
	if err != nil {
		if err == ErrExpiredKey {
			tx.db.evict(key, String)
		}
		return false
	}

	return true
}

func (tx *Tx) set(key string, value string) error {
	var existVal string
	existVal, err := tx.get(key)
	if err != nil && err != ErrExpiredKey && err != ErrInvalidKey {
		return err
	}

	if reflect.DeepEqual(existVal, value) {
		return err
	}

	tx.db.strStore.Set(key, value)

	return nil
}

func (tx *Tx) get(key string) (val string, err error) {
	v, err := tx.db.strStore.get(key)
	if err != nil {
		return "", err
	}

	// Check if the key is expired.
	if tx.db.hasExpired(key, String) {
		tx.db.evict(key, String)
		return "", ErrExpiredKey
	}

	val = v.(string)
	return
}
