package flashdb

import (
	"os"
	"testing"

	"github.com/arriqaaq/aol"
	"github.com/stretchr/testify/assert"
)

func TestRoseDB_ZSet(t *testing.T) {
	db := getTestDB()
	logPath := "tmp/"
	l, err := aol.Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	db.log = l
	defer os.RemoveAll("tmp/")

	err = db.ZAdd(testKey, 1, "foo")
	assert.NoError(t, err)
	err = db.ZAdd(testKey, 2, "bar")
	assert.NoError(t, err)
	err = db.ZAdd(testKey, 3, "baz")
	assert.NoError(t, err)

	_, s := db.ZScore(testKey, "foo")
	assert.Equal(t, 1.0, s)

	assert.Equal(t, 3, db.ZCard(testKey))
}
