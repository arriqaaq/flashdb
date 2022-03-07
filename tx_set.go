package flashdb

import (
	"time"
)

func (tx *Tx) SAdd(key string, members ...string) (err error) {
	for _, m := range members {
		exist := tx.db.setStore.SIsMember(key, m)
		if !exist {
			e := newRecord([]byte(key), []byte(m), SetRecord, SetSAdd)
			tx.addRecord(e)
		}
	}
	return
}

func (tx *Tx) SIsMember(key string, member string) bool {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return false
	}
	return tx.db.setStore.SIsMember(key, member)
}

func (tx *Tx) SRandMember(key string, count int) (values []string) {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return nil
	}

	vals := tx.db.setStore.SRandMember(key, count)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (tx *Tx) SRem(key string, members ...string) (res int, err error) {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return
	}

	for _, m := range members {
		e := newRecord([]byte(key), []byte(m), SetRecord, SetSRem)
		tx.addRecord(e)
		res++
	}
	return
}

func (tx *Tx) SMove(src, dst string, member string) error {
	if tx.db.hasExpired(src, Set) {
		tx.db.evict(src, Hash)
		return ErrExpiredKey
	}
	if tx.db.hasExpired(dst, Set) {
		tx.db.evict(dst, Hash)
		return ErrExpiredKey
	}

	ok := tx.db.setStore.SMove(src, dst, member)
	if ok {
		e := newRecordWithValue([]byte(src), []byte(member), []byte(dst), SetRecord, SetSMove)
		tx.addRecord(e)
	}
	return nil
}

func (tx *Tx) SCard(key string) int {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return 0
	}
	return tx.db.setStore.SCard(key)
}

func (tx *Tx) SMembers(key string) (values []string) {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return
	}

	vals := tx.db.setStore.SMembers(key)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (tx *Tx) SUnion(keys ...string) (values []string) {
	var activeKeys []string
	for _, k := range keys {
		if tx.db.hasExpired(k, Set) {
			tx.db.evict(k, Hash)
			continue
		}
		activeKeys = append(activeKeys, k)
	}

	vals := tx.db.setStore.SUnion(activeKeys...)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

func (tx *Tx) SDiff(keys ...string) (values []string) {
	var activeKeys []string
	for _, k := range keys {
		if tx.db.hasExpired(k, Set) {
			tx.db.evict(k, Hash)
			continue
		}
		activeKeys = append(activeKeys, k)
	}

	vals := tx.db.setStore.SDiff(activeKeys...)
	for _, v := range vals {
		values = append(values, toString(v))
	}
	return
}

// SKeyExists returns if the key exists.
func (tx *Tx) SKeyExists(key string) (ok bool) {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)

		return
	}

	ok = tx.db.setStore.SKeyExists(key)
	return
}

// SClear clear the specified key in set.
func (tx *Tx) SClear(key string) (err error) {
	if !tx.SKeyExists(key) {
		return ErrInvalidKey
	}

	e := newRecord([]byte(key), nil, SetRecord, SetSClear)
	tx.addRecord(e)
	return
}

// SExpire set expired time for the key in set.
func (tx *Tx) SExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}
	if !tx.SKeyExists(key) {
		return ErrInvalidKey
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, SetRecord, SetSExpire)
	tx.addRecord(e)
	return
}

// STTL return time to live for the key in set.
func (tx *Tx) STTL(key string) (ttl int64) {
	if tx.db.hasExpired(key, Set) {
		tx.db.evict(key, Set)
		return
	}

	deadline := tx.db.getTTL(Set, key)
	if deadline == nil {
		return
	}

	return deadline.(int64) - time.Now().Unix()
}
