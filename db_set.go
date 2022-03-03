package flashdb

import (
	"time"
)

func (db *FlashDB) SAdd(key string, members ...string) (res int, err error) {
	db.setStore.Lock()
	defer db.setStore.Unlock()

	for _, m := range members {
		exist := db.setStore.SIsMember(key, m)
		if !exist {
			e := newRecord([]byte(key), []byte(m), SetRecord, SetSAdd)
			if err = db.write(e); err != nil {
				return
			}
			res = db.setStore.SAdd(key, m)
		}
	}
	return
}

func (db *FlashDB) SPop(key string, count int) (values []string, err error) {
	db.setStore.Lock()
	defer db.setStore.Unlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return nil, ErrExpiredKey
	}

	vals := db.setStore.SPop(key, count)
	for _, v := range vals {
		e := newRecord([]byte(key), []byte(toString(v)), SetRecord, SetSRem)
		if err = db.write(e); err != nil {
			return
		}
		values = append(values, toString(v))
	}
	return
}

func (db *FlashDB) SIsMember(key string, member string) bool {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return false
	}
	return db.setStore.SIsMember(key, member)
}

func (db *FlashDB) SRandMember(key string, count int) (values []string) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return nil
	}

	vals := db.setStore.SRandMember(key, count)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (db *FlashDB) SRem(key string, members ...string) (res int, err error) {
	db.setStore.Lock()
	defer db.setStore.Unlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return
	}

	for _, m := range members {
		e := newRecord([]byte(key), []byte(m), SetRecord, SetSRem)
		if err = db.write(e); err != nil {
			return
		}
		if ok := db.setStore.SRem(key, m); ok {
			res++
		}
	}
	return
}

func (db *FlashDB) SMove(src, dst string, member string) error {
	db.setStore.Lock()
	defer db.setStore.Unlock()

	if db.hasExpired(src, Set) {
		db.evict(src, Hash)
		return ErrExpiredKey
	}
	if db.hasExpired(dst, Set) {
		db.evict(dst, Hash)
		return ErrExpiredKey
	}

	ok := db.setStore.SMove(src, dst, member)
	if ok {
		e := newRecordWithValue([]byte(src), []byte(member), []byte(dst), SetRecord, SetSMove)
		if err := db.write(e); err != nil {
			return err
		}
	}
	return nil
}

func (db *FlashDB) SCard(key string) int {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return 0
	}
	return db.setStore.SCard(key)
}

func (db *FlashDB) SMembers(key string) (values []string) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return
	}

	vals := db.setStore.SMembers(key)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (db *FlashDB) SUnion(keys ...string) (values []string) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	var activeKeys []string
	for _, k := range keys {
		if db.hasExpired(k, Set) {
			db.evict(k, Hash)
			continue
		}
		activeKeys = append(activeKeys, k)
	}

	vals := db.setStore.SUnion(activeKeys...)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (db *FlashDB) SDiff(keys ...string) (values []string) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	var activeKeys []string
	for _, k := range keys {
		if db.hasExpired(k, Set) {
			db.evict(k, Hash)
			continue
		}
		activeKeys = append(activeKeys, k)
	}

	vals := db.setStore.SDiff(activeKeys...)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

// SKeyExists returns if the key exists.
func (db *FlashDB) SKeyExists(key string) (ok bool) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)

		return
	}

	ok = db.setStore.SKeyExists(key)
	return
}

// SClear clear the specified key in set.
func (db *FlashDB) SClear(key string) (err error) {
	if !db.SKeyExists(key) {
		return ErrInvalidKey
	}

	db.setStore.Lock()
	defer db.setStore.Unlock()

	e := newRecord([]byte(key), nil, SetRecord, SetSClear)
	if err = db.write(e); err != nil {
		return
	}
	db.setStore.SClear(key)
	db.exps.HDel(Set, key)
	return
}

// SExpire set expired time for the key in set.
func (db *FlashDB) SExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}
	if !db.SKeyExists(key) {
		return ErrInvalidKey
	}

	db.setStore.Lock()
	defer db.setStore.Unlock()

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, SetRecord, SetSExpire)
	if err = db.write(e); err != nil {
		return
	}

	db.setTTL(Set, key, ttl)
	return
}

// STTL return time to live for the key in set.
func (db *FlashDB) STTL(key string) (ttl int64) {
	db.setStore.RLock()
	defer db.setStore.RUnlock()

	if db.hasExpired(key, Set) {
		db.evict(key, Set)
		return
	}

	deadline := db.getTTL(Set, key)
	if deadline == nil {
		return
	}

	return deadline.(int64) - time.Now().Unix()
}
