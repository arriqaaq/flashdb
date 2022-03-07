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
	db.persist = true
	defer os.RemoveAll("tmp/")

	if err := db.Update(func(tx *Tx) error {
		err = tx.ZAdd(testKey, 1, "foo")
		assert.NoError(t, err)
		err = tx.ZAdd(testKey, 2, "bar")
		assert.NoError(t, err)
		err = tx.ZAdd(testKey, 3, "baz")
		assert.NoError(t, err)

		_, s := tx.ZScore(testKey, "foo")
		assert.Equal(t, 1.0, s)

		assert.Equal(t, 3, tx.ZCard(testKey))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
