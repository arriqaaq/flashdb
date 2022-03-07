package flashdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlashDB_GetSet(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		err := tx.Set("foo", "bar")
		assert.Equal(t, nil, err)
		val, err := tx.Get("foo")
		assert.Equal(t, nil, err)
		assert.Equal(t, "bar", val)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_SetEx(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		err := tx.SetEx("foo", "1", -4)
		assert.NotEmpty(t, err)

		err = tx.SetEx("foo", "1", 993)
		assert.Empty(t, err)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_Delete(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		err := tx.Set("foo", "bar")
		assert.Equal(t, err, nil)

		err = tx.Delete("foo")
		assert.Equal(t, err, nil)

		val, err := tx.Get("foo")
		assert.Empty(t, val)
		assert.Equal(t, ErrInvalidKey, err)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFlashDB_TTL(t *testing.T) {
	db := getTestDB()
	if err := db.View(func(tx *Tx) error {
		err := tx.SetEx("foo", "bar", 20)
		assert.Equal(t, err, nil)

		ttl := tx.TTL("foo")
		assert.Equal(t, int(ttl), 20)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
