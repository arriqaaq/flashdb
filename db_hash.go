package flashdb

import (
	"reflect"
	"time"
)

func (db *FlashDB) HSet(key string, field string, value string) (res int, err error) {
	existVal := db.HGet(key, field)
	if reflect.DeepEqual(existVal, value) {
		return
	}

	db.hashStore.Lock()
	defer db.hashStore.Unlock()

	e := newRecordWithValue([]byte(key), []byte(field), []byte(value), HashRecord, HashHSet)
	if err = db.write(e); err != nil {
		return
	}

	res = db.hashStore.HSet(key, field, value)
	return
}

func (db *FlashDB) HGet(key string, field string) string {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return ""
	}

	return toString(db.hashStore.HGet(key, field))
}

func (db *FlashDB) HGetAll(key string) []string {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return nil
	}

	vals := db.hashStore.HGetAll(key)
	values := make([]string, 0, 1)

	for _, v := range vals {
		values = append(values, toString(v))
	}

	return values
}

func (db *FlashDB) HDel(key string, field ...string) (res int, err error) {
	db.hashStore.Lock()
	defer db.hashStore.Unlock()

	for _, f := range field {
		if ok := db.hashStore.HDel(key, f); ok == 1 {
			e := newRecord([]byte(key), nil, HashRecord, HashHDel)
			if err = db.write(e); err != nil {
				return
			}
			res++
		}
	}
	return
}

func (db *FlashDB) HKeyExists(key string) (ok bool) {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return
	}
	return db.hashStore.HKeyExists(key)
}

func (db *FlashDB) HExists(key, field string) (ok bool) {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return
	}

	return db.hashStore.HExists(key, field)
}

func (db *FlashDB) HLen(key string) int {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return 0
	}

	return db.hashStore.HLen(key)
}

func (db *FlashDB) HKeys(key string) (val []string) {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return nil
	}

	return db.hashStore.HKeys(key)
}

func (db *FlashDB) HVals(key string) (values []string) {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return nil
	}

	vals := db.hashStore.HVals(key)
	for _, v := range vals {
		values = append(values, toString(v))
	}

	return
}

func (db *FlashDB) HExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	if !db.HKeyExists(key) {
		return ErrInvalidKey
	}

	db.hashStore.Lock()
	defer db.hashStore.Unlock()

	ttl := time.Now().Unix() + duration

	db.setTTL(Hash, key, ttl)
	return
}

func (db *FlashDB) HTTL(key string) (ttl int64) {
	db.hashStore.RLock()
	defer db.hashStore.RUnlock()

	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return
	}

	deadline := db.getTTL(Hash, key)
	if deadline == nil {
		return
	}
	return deadline.(int64) - time.Now().Unix()
}

func (db *FlashDB) HClear(key string) (err error) {
	if db.hasExpired(key, Hash) {
		db.evict(key, Hash)
		return
	}

	db.hashStore.Lock()
	defer db.hashStore.Unlock()

	e := newRecord([]byte(key), nil, HashRecord, HashHClear)
	if err := db.write(e); err != nil {
		return err
	}

	db.hashStore.HClear(key)
	db.exps.HDel(Hash, key)
	return
}

func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}
