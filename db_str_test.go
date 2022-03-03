package flashdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestDB() *FlashDB {
	db, _ := New(DefaultConfig())
	return db
}

func TestFlashDB_GetSet(t *testing.T) {
	db := getTestDB()
	err := db.Set("foo", "bar")
	assert.Equal(t, nil, err)
	val, err := db.Get("foo")
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", val)
}

func TestFlashDB_SetEx(t *testing.T) {
	db := getTestDB()
	err := db.SetEx("foo", "1", -4)
	assert.NotEmpty(t, err)

	err = db.SetEx("foo", "1", 993)
	assert.Empty(t, err)
}

func TestFlashDB_Delete(t *testing.T) {
	db := getTestDB()
	err := db.Set("foo", "bar")
	assert.Equal(t, err, nil)

	err = db.Delete("foo")
	assert.Equal(t, err, nil)

	val, err := db.Get("foo")
	assert.Empty(t, val)
	assert.Equal(t, ErrInvalidKey, err)
}

func TestFlashDB_TTL(t *testing.T) {
	db := getTestDB()

	err := db.SetEx("foo", "bar", 20)
	assert.Equal(t, err, nil)

	ttl := db.TTL("foo")
	assert.Equal(t, int(ttl), 20)
}
