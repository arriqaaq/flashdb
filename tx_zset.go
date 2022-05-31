package flashdb

import (
	"time"
)

// ZAdd adds key-member pair with the score. If the key-member pair already
// exists and the old score is the same as the new score, it doesn't do anything.
func (tx *Tx) ZAdd(key string, score float64, member string) error {
	if ok, oldScore := tx.ZScore(key, member); ok && oldScore == score {
		return nil
	}

	value := float64ToStr(score)
	e := newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd)
	tx.addRecord(e)

	return nil
}

// ZScore returns score of the given key-member pair.If the key has expired,
// the key is evicted.
func (tx *Tx) ZScore(key string, member string) (ok bool, score float64) {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return
	}

	return tx.db.zsetStore.ZScore(key, member)
}

// ZCard returns sorted set cardinality(number of elements) of the sorted set
// stored at key. If the key has expired, the key is evicted.
func (tx *Tx) ZCard(key string) int {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return 0
	}

	return tx.db.zsetStore.ZCard(key)
}

// ZRank returns the rank of the member at key, with the scores ordered from
// low to high. If the key has expired, the key is evicted.
func (tx *Tx) ZRank(key string, member string) int64 {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return -1
	}

	return tx.db.zsetStore.ZRank(key, member)
}

// ZRevRank returns the rank of the member at key, with the scores ordered from
// high to low. If the key has expired, the key is evicted.
func (tx *Tx) ZRevRank(key string, member string) int64 {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return -1
	}

	return tx.db.zsetStore.ZRevRank(key, member)
}

// ZRange returns the specified range of elements in the sorted set stored at
// key. If the key has expired, the key is evicted.
func (tx *Tx) ZRange(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRange(key, start, stop)
}

// ZRangeWithScores returns the specified range of elements with scores in the
// sorted set stored at key. If the key has expired, the key is evicted.
func (tx *Tx) ZRangeWithScores(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRangeWithScores(key, start, stop)
}

// ZRevRange returns the specified range of elements in the sorted set stored at
// key. The elements are ordered from the highest score to the lowest score. If
// key has expired, the key is evicted.
func (tx *Tx) ZRevRange(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevRange(key, start, stop)
}

// ZRevRangeWithScores returns the specified range of elements in the sorted set
// at key. The elements are ordered from the highest to the lowest score. If key
// has expired, the key is evicted.
func (tx *Tx) ZRevRangeWithScores(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevRangeWithScores(key, start, stop)
}

// ZRem removes the member from the sorted set at key.
func (tx *Tx) ZRem(key string, member string) (ok bool, err error) {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return
	}

	ok = tx.db.zsetStore.ZRem(key, member)
	if ok {
		e := newRecord([]byte(key), []byte(member), ZSetRecord, ZSetZRem)
		tx.addRecord(e)
	}

	return
}

// ZGetByRank returns the members by given rank at key. If the key has expired,
// the key is evicted.
func (tx *Tx) ZGetByRank(key string, rank int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZGetByRank(key, rank)
}

// ZRevGetByRank returns the members by given rank at key. The members are
// returned reverse ordered. If the key has expired, the key is evicted.
func (tx *Tx) ZRevGetByRank(key string, rank int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevGetByRank(key, rank)
}

// ZScoreRange returns the members in given range at key. If the key has expired,
// the key is evicted.
func (tx *Tx) ZScoreRange(key string, min, max float64) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZScoreRange(key, min, max)
}

// ZRevScoreRange returns the members in given range at key. The members are
// returned in reverse order. If the key has expired, the key is evicted.
func (tx *Tx) ZRevScoreRange(key string, max, min float64) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevScoreRange(key, max, min)
}

// ZKeyExists checks the sorted set whether the key exists. If the key has expired,
// the key is evicted.
func (tx *Tx) ZKeyExists(key string) (ok bool) {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return
	}

	ok = tx.db.zsetStore.ZKeyExists(key)
	return
}

// ZClear clears the members at key.
func (tx *Tx) ZClear(key string) (err error) {
	e := newRecord([]byte(key), nil, ZSetRecord, ZSetZClear)
	tx.addRecord(e)

	return
}

// ZExpire sets expire time at key. duration should be more than zero.
func (tx *Tx) ZExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}
	if !tx.ZKeyExists(key) {
		return ErrInvalidKey
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, ZSetRecord, ZSetZExpire)
	tx.addRecord(e)
	return
}

// ZTTL returns the remaining TTL of the given key.
func (tx *Tx) ZTTL(key string) (ttl int64) {
	if !tx.ZKeyExists(key) {
		return
	}

	deadline := tx.db.getTTL(ZSet, key)
	if deadline == nil {
		return
	}
	return deadline.(int64) - time.Now().Unix()
}
