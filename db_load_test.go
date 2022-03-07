package flashdb

import (
	"fmt"
	"os"
	"testing"

	"github.com/arriqaaq/aol"
	"github.com/stretchr/testify/assert"
)

func makeLoadRecords(n int, db *FlashDB) {
	for i := 1; i <= n; i++ {
		key := fmt.Sprintf("key_%d", i)
		member := fmt.Sprintf("member_%d", i)
		value := fmt.Sprintf("value_%d", i)
		db.Update(func(tx *Tx) error {
			tx.HSet(key, member, value)
			tx.SAdd(key, member)
			tx.Set(key, member)
			tx.ZAdd(key, 10.0, member)
			return nil
		})

	}
}

func TestFlashDB_load(t *testing.T) {
	db := getTestDB()
	logPath := "tmp/"
	l, err := aol.Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	db.log = l
	db.persist = true
	defer os.RemoveAll("tmp/")

	makeLoadRecords(10, db)

	db.Close()

	p, err := aol.Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	db2 := getTestDB()
	db2.log = p
	err = db2.load()
	assert.NoError(t, err)

	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		member := fmt.Sprintf("member_%d", i)
		value := fmt.Sprintf("value_%d", i)
		db2.View(func(tx *Tx) error {
			assert.Equal(t, value, tx.HGet(key, member))
			assert.True(t, tx.SIsMember(key, member))
			val, err := tx.Get(key)
			assert.NoError(t, err)
			assert.Equal(t, member, val)
			ok, score := tx.ZScore(key, member)
			assert.True(t, ok)
			assert.Equal(t, 10.0, score)
			return nil
		})

	}
}
