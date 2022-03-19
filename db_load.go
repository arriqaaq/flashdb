package flashdb

import (
	"time"

	"github.com/arriqaaq/aol"
)

// load String, Hash, Set and ZSet stores from append-only log
func (db *FlashDB) load() error {
	if db.log == nil {
		return nil
	}

	noOfSegments := db.log.Segments()
	for i := 1; i <= noOfSegments; i++ {
		j := 0

		for {
			data, err := db.log.Read(uint64(i), uint64(j))
			if err != nil {
				if err == aol.ErrEOF {
					break
				}
				return err
			}

			record, err := decode(data)
			if err != nil {
				return err
			}

			if len(record.meta.key) > 0 {
				if err := db.loadRecord(record); err != nil {
					return err
				}
			}

			j++
		}
	}

	return nil
}

func (db *FlashDB) loadRecord(r *record) (err error) {

	switch r.getType() {
	case StringRecord:
		err = db.buildStringRecord(r)
	case HashRecord:
		err = db.buildHashRecord(r)
	case SetRecord:
		err = db.buildSetRecord(r)
	case ZSetRecord:
		err = db.buildZsetRecord(r)
	}
	return
}

/*
	Utility functions to build stores from aol Record
*/

func (db *FlashDB) buildStringRecord(r *record) error {

	key := string(r.meta.key)
	member := string(r.meta.member)

	switch r.getMark() {
	case StringSet:
		db.strStore.Insert([]byte(key), member)
	case StringRem:
		db.strStore.Delete([]byte(key))
		db.exps.HDel(String, key)
	case StringExpire:
		if r.timestamp < uint64(time.Now().Unix()) {
			db.strStore.Delete([]byte(key))
			db.exps.HDel(String, key)
		} else {
			db.setTTL(String, key, int64(r.timestamp))
		}
	}

	return nil
}

func (db *FlashDB) buildHashRecord(r *record) error {

	key := string(r.meta.key)
	member := string(r.meta.member)
	value := string(r.meta.value)

	switch r.getMark() {
	case HashHSet:
		db.hashStore.HSet(key, member, value)
	case HashHDel:
		db.hashStore.HDel(key, member)
	case HashHClear:
		db.hashStore.HClear(key)
		db.exps.HDel(Hash, key)
	case HashHExpire:
		if r.timestamp < uint64(time.Now().Unix()) {
			db.hashStore.HClear(key)
			db.exps.HDel(Hash, key)
		} else {
			db.setTTL(Hash, key, int64(r.timestamp))
		}
	}

	return nil
}

func (db *FlashDB) buildSetRecord(r *record) error {

	key := string(r.meta.key)
	member := string(r.meta.member)
	value := string(r.meta.value)

	switch r.getMark() {
	case SetSAdd:
		db.setStore.SAdd(key, member)
	case SetSRem:
		db.setStore.SRem(key, member)
	case SetSMove:
		db.setStore.SMove(key, value, member)
	case SetSClear:
		db.setStore.SClear(key)
		db.exps.HDel(Set, key)
	case SetSExpire:
		if r.timestamp < uint64(time.Now().Unix()) {
			db.setStore.SClear(key)
			db.exps.HDel(Set, key)
		} else {
			db.setTTL(Set, key, int64(r.timestamp))
		}
	}

	return nil
}

func (db *FlashDB) buildZsetRecord(r *record) error {

	key := string(r.meta.key)
	member := string(r.meta.member)
	value := string(r.meta.value)

	switch r.getMark() {
	case ZSetZAdd:
		score, err := strToFloat64(value)
		if err != nil {
			return err
		}
		db.zsetStore.ZAdd(key, score, member, nil)
	case ZSetZRem:
		db.zsetStore.ZRem(key, member)
	case ZSetZClear:
		db.zsetStore.ZClear(key)
		db.exps.HDel(ZSet, key)
	case ZSetZExpire:
		if r.timestamp < uint64(time.Now().Unix()) {
			db.zsetStore.ZClear(key)
			db.exps.HDel(ZSet, key)
		} else {
			db.setTTL(ZSet, key, int64(r.timestamp))
		}
	}

	return nil
}
