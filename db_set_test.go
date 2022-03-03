package flashdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlashDB_SAddPop(t *testing.T) {
	db := getTestDB()
	db.SAdd(testKey, "foo", "bar", "baz")
	values, err := db.SPop(testKey, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(values))
}

func TestFlashDB_SCard(t *testing.T) {
	db := getTestDB()
	db.SAdd(testKey, "foo", "bar", "baz")

	cnt := db.SCard(testKey)
	assert.Equal(t, 3, cnt)
}

func TestFlashDB_SIsMember(t *testing.T) {
	db := getTestDB()
	db.SAdd(testKey, "foo", "bar", "baz")
	assert.True(t, db.SIsMember(testKey, "foo"))
	assert.True(t, db.SIsMember(testKey, "bar"))
	assert.True(t, db.SIsMember(testKey, "baz"))
}

func TestFlashDB_SRem(t *testing.T) {
	db := getTestDB()
	db.SAdd(testKey, "foo", "bar", "baz")
	db.SRem(testKey, "foo")
	assert.False(t, db.SIsMember(testKey, "foo"))
}

func TestFlashDB_SClear(t *testing.T) {
	db := getTestDB()
	db.SAdd(testKey, "foo", "bar", "baz")
	assert.NoError(t, db.SClear(testKey))
	assert.False(t, db.SKeyExists(testKey))
}

func TestFlashDB_SDiff(t *testing.T) {
	db := getTestDB()
	db.SAdd("set1", "foo", "bar", "baz")
	db.SAdd("set2", "foo", "bar")
	res := db.SDiff("set1", "set2")
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "baz", res[0])
}
