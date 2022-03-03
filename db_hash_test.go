package flashdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKey = "dummy"

func TestFlashDB_HGetSet(t *testing.T) {

	db := getTestDB()
	db.HSet(testKey, "bar", "1")
	db.HSet(testKey, "baz", "2")

	val := db.HGet(testKey, "bar")
	assert.Equal(t, "1", val)
	val = db.HGet(testKey, "baz")
	assert.Equal(t, "2", val)

}

func TestFlashDB_HGetAll(t *testing.T) {
	db := getTestDB()
	db.HSet(testKey, "bar", "1")
	db.HSet(testKey, "baz", "2")

	values := db.HGetAll(testKey)
	assert.Equal(t, 4, len(values))
}

func TestFlashDB_HDel(t *testing.T) {
	db := getTestDB()
	db.HSet(testKey, "bar", "1")
	db.HSet(testKey, "baz", "2")

	res, err := db.HDel(testKey, "bar", "baz")
	assert.Nil(t, err)
	assert.Equal(t, 2, res)
	assert.Empty(t, db.HGet(testKey, "bar"))
	assert.Empty(t, db.HGet(testKey, "baz"))
}

func TestFlashDB_HExists(t *testing.T) {
	db := getTestDB()
	db.HSet(testKey, "bar", "1")
	db.HSet(testKey, "baz", "2")

	assert.True(t, db.HExists(testKey, "bar"))
	assert.True(t, db.HExists(testKey, "baz"))
	assert.False(t, db.HExists(testKey, "ben"))
}
