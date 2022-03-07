package flashdb

import (
	"time"
)

func (tx *Tx) ZAdd(key string, score float64, member string) error {
	if ok, oldScore := tx.ZScore(key, member); ok && oldScore == score {
		return nil
	}

	value := float64ToStr(score)
	e := newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd)
	tx.addRecord(e)

	tx.db.zsetStore.ZAdd(key, score, member, nil)
	return nil
}

func (tx *Tx) ZScore(key string, member string) (ok bool, score float64) {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return
	}

	return tx.db.zsetStore.ZScore(key, member)
}

func (tx *Tx) ZCard(key string) int {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return 0
	}

	return tx.db.zsetStore.ZCard(key)
}

func (tx *Tx) ZRank(key string, member string) int64 {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return -1
	}

	return tx.db.zsetStore.ZRank(key, member)
}

func (tx *Tx) ZRevRank(key string, member string) int64 {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return -1
	}

	return tx.db.zsetStore.ZRevRank(key, member)
}

func (tx *Tx) ZIncrBy(key string, increment float64, member string) (float64, error) {
	value := float64ToStr(increment)
	e := newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd)
	tx.addRecord(e)

	increment = tx.db.zsetStore.ZIncrBy(key, increment, member)

	return increment, nil
}

func (tx *Tx) ZRange(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRange(key, start, stop)
}

func (tx *Tx) ZRangeWithScores(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRangeWithScores(key, start, stop)
}

func (tx *Tx) ZRevRange(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevRange(key, start, stop)
}

func (tx *Tx) ZRevRangeWithScores(key string, start, stop int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevRangeWithScores(key, start, stop)
}

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

func (tx *Tx) ZGetByRank(key string, rank int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZGetByRank(key, rank)
}

func (tx *Tx) ZRevGetByRank(key string, rank int) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevGetByRank(key, rank)
}

func (tx *Tx) ZScoreRange(key string, min, max float64) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZScoreRange(key, min, max)
}

func (tx *Tx) ZRevScoreRange(key string, max, min float64) []interface{} {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return nil
	}

	return tx.db.zsetStore.ZRevScoreRange(key, max, min)
}

func (tx *Tx) ZKeyExists(key string) (ok bool) {
	if tx.db.hasExpired(key, ZSet) {
		tx.db.evict(key, ZSet)
		return
	}

	ok = tx.db.zsetStore.ZKeyExists(key)
	return
}

func (tx *Tx) ZClear(key string) (err error) {
	e := newRecord([]byte(key), nil, ZSetRecord, ZSetZClear)
	tx.addRecord(e)

	tx.db.zsetStore.ZClear(key)
	tx.db.exps.HDel(ZSet, key)
	return
}

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

	tx.db.setTTL(ZSet, key, ttl)
	return
}

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
