package flashdb

import (
	"reflect"
	"time"
)

func (tx *Tx) HSet(key string, field string, value string) (res int, err error) {
	existVal := tx.HGet(key, field)
	if reflect.DeepEqual(existVal, value) {
		return
	}

	e := newRecordWithValue([]byte(key), []byte(field), []byte(value), HashRecord, HashHSet)
	tx.addRecord(e)

	res = tx.db.hashStore.HSet(key, field, value)
	return
}

func (tx *Tx) HGet(key string, field string) string {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return ""
	}

	return toString(tx.db.hashStore.HGet(key, field))
}

func (tx *Tx) HGetAll(key string) []string {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	vals := tx.db.hashStore.HGetAll(key)
	values := make([]string, 0, 1)

	for _, v := range vals {
		values = append(values, toString(v))
	}

	return values
}

func (tx *Tx) HDel(key string, field ...string) (res int, err error) {

	for _, f := range field {
		if ok := tx.db.hashStore.HDel(key, f); ok == 1 {
			e := newRecord([]byte(key), nil, HashRecord, HashHDel)
			if tx.db.persist {
				tx.wc.rollbackItems = append(tx.wc.rollbackItems, e)
				tx.wc.commitItems = append(tx.wc.commitItems, e)
			}
			res++
		}
	}
	return
}

func (tx *Tx) HKeyExists(key string) (ok bool) {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}
	return tx.db.hashStore.HKeyExists(key)
}

func (tx *Tx) HExists(key, field string) (ok bool) {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	return tx.db.hashStore.HExists(key, field)
}

func (tx *Tx) HLen(key string) int {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return 0
	}

	return tx.db.hashStore.HLen(key)
}

func (tx *Tx) HKeys(key string) (val []string) {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	return tx.db.hashStore.HKeys(key)
}

func (tx *Tx) HVals(key string) (values []string) {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	vals := tx.db.hashStore.HVals(key)
	for _, v := range vals {
		values = append(values, toString(v))
	}

	return
}

func (tx *Tx) HExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	if !tx.HKeyExists(key) {
		return ErrInvalidKey
	}

	ttl := time.Now().Unix() + duration

	tx.db.setTTL(Hash, key, ttl)
	return
}

func (tx *Tx) HTTL(key string) (ttl int64) {
	tx.db.hashStore.RLock()
	defer tx.db.hashStore.RUnlock()

	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	deadline := tx.db.getTTL(Hash, key)
	if deadline == nil {
		return
	}
	return deadline.(int64) - time.Now().Unix()
}

func (tx *Tx) HClear(key string) (err error) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	e := newRecord([]byte(key), nil, HashRecord, HashHClear)
	tx.addRecord(e)

	tx.db.hashStore.HClear(key)
	tx.db.exps.HDel(Hash, key)
	return
}

func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}
