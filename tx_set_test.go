package flashdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlashDB_SAddPop(t *testing.T) {
	db := getTestDB()
	if err := db.Update(func(tx *Tx) error {
		tx.SAdd(testKey, "foo", "bar", "baz")
		values, err := tx.SPop(testKey, 2)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(values))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_SCard(t *testing.T) {
	db := getTestDB()
	if err := db.Update(func(tx *Tx) error {
		tx.SAdd(testKey, "foo", "bar", "baz")
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	db.View(func(tx *Tx) error {
		cnt := tx.SCard(testKey)
		assert.Equal(t, 3, cnt)
		return nil
	})
}

func TestFlashDB_SIsMember(t *testing.T) {
	db := getTestDB()
	if err := db.Update(func(tx *Tx) error {
		tx.SAdd(testKey, "foo", "bar", "baz")
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	db.View(func(tx *Tx) error {
		assert.True(t, tx.SIsMember(testKey, "foo"))
		assert.True(t, tx.SIsMember(testKey, "bar"))
		assert.True(t, tx.SIsMember(testKey, "baz"))
		return nil
	})
}

func TestFlashDB_SRem(t *testing.T) {
	db := getTestDB()
	if err := db.Update(func(tx *Tx) error {
		tx.SAdd(testKey, "foo", "bar", "baz")
		tx.SRem(testKey, "foo")
		assert.False(t, tx.SIsMember(testKey, "foo"))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_SClear(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		tx.SAdd(testKey, "foo", "bar", "baz")
		assert.NoError(t, tx.SClear(testKey))
		assert.False(t, tx.SKeyExists(testKey))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_SDiff(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		tx.SAdd("set1", "foo", "bar", "baz")
		tx.SAdd("set2", "foo", "bar")
		res := tx.SDiff("set1", "set2")
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "baz", res[0])
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
