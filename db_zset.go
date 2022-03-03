package flashdb

import (
	"time"
)

func (db *FlashDB) ZAdd(key string, score float64, member string) error {
	if ok, oldScore := db.ZScore(key, member); ok && oldScore == score {
		return nil
	}

	db.zsetStore.Lock()
	defer db.zsetStore.Unlock()

	value := float64ToStr(score)
	e := newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd)
	if err := db.write(e); err != nil {
		return err
	}

	db.zsetStore.ZAdd(key, score, member, nil)
	return nil
}

func (db *FlashDB) ZScore(key string, member string) (ok bool, score float64) {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return
	}

	return db.zsetStore.ZScore(key, member)
}

func (db *FlashDB) ZCard(key string) int {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return 0
	}

	return db.zsetStore.ZCard(key)
}

func (db *FlashDB) ZRank(key string, member string) int64 {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return -1
	}

	return db.zsetStore.ZRank(key, member)
}

func (db *FlashDB) ZRevRank(key string, member string) int64 {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return -1
	}

	return db.zsetStore.ZRevRank(key, member)
}

func (db *FlashDB) ZIncrBy(key string, increment float64, member string) (float64, error) {
	db.zsetStore.Lock()
	defer db.zsetStore.Unlock()

	value := float64ToStr(increment)
	e := newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd)
	if err := db.write(e); err != nil {
		return increment, err
	}

	increment = db.zsetStore.ZIncrBy(key, increment, member)

	return increment, nil
}

func (db *FlashDB) ZRange(key string, start, stop int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRange(key, start, stop)
}

func (db *FlashDB) ZRangeWithScores(key string, start, stop int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRangeWithScores(key, start, stop)
}

func (db *FlashDB) ZRevRange(key string, start, stop int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRevRange(key, start, stop)
}

func (db *FlashDB) ZRevRangeWithScores(key string, start, stop int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRevRangeWithScores(key, start, stop)
}

func (db *FlashDB) ZRem(key string, member string) (ok bool, err error) {
	db.zsetStore.Lock()
	defer db.zsetStore.Unlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return
	}

	ok = db.zsetStore.ZRem(key, member)
	if ok {
		e := newRecord([]byte(key), []byte(member), ZSetRecord, ZSetZRem)
		if err := db.write(e); err != nil {
			return ok, err
		}
	}

	return
}

func (db *FlashDB) ZGetByRank(key string, rank int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZGetByRank(key, rank)
}

func (db *FlashDB) ZRevGetByRank(key string, rank int) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRevGetByRank(key, rank)
}

func (db *FlashDB) ZScoreRange(key string, min, max float64) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZScoreRange(key, min, max)
}

func (db *FlashDB) ZRevScoreRange(key string, max, min float64) []interface{} {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return nil
	}

	return db.zsetStore.ZRevScoreRange(key, max, min)
}

func (db *FlashDB) ZKeyExists(key string) (ok bool) {
	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	if db.hasExpired(key, ZSet) {
		db.evict(key, ZSet)
		return
	}

	ok = db.zsetStore.ZKeyExists(key)
	return
}

func (db *FlashDB) ZClear(key string) (err error) {
	db.zsetStore.Lock()
	defer db.zsetStore.Unlock()

	e := newRecord([]byte(key), nil, ZSetRecord, ZSetZClear)
	if err := db.write(e); err != nil {
		return err
	}

	db.zsetStore.ZClear(key)
	db.exps.HDel(ZSet, key)
	return
}

func (db *FlashDB) ZExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}
	if !db.ZKeyExists(key) {
		return ErrInvalidKey
	}

	db.zsetStore.Lock()
	defer db.zsetStore.Unlock()

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, ZSetRecord, ZSetZExpire)
	if err := db.write(e); err != nil {
		return err
	}

	db.setTTL(ZSet, key, ttl)
	return
}

func (db *FlashDB) ZTTL(key string) (ttl int64) {
	if !db.ZKeyExists(key) {
		return
	}

	db.zsetStore.RLock()
	defer db.zsetStore.RUnlock()

	deadline := db.getTTL(ZSet, key)
	if deadline == nil {
		return
	}
	return deadline.(int64) - time.Now().Unix()
}
