package flashdb

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func newKey(n int) string {
	return "test_key_" + fmt.Sprintf("%09d", n)
}

func newValue(n int) string {
	return "test_val-" + fmt.Sprintf("%09d", n)
}

func randomValue() string {
	var str bytes.Buffer
	for i := 0; i < 12; i++ {
		str.WriteByte(alphabet[rand.Int()%26])
	}
	return "test_val-" + strconv.FormatInt(time.Now().UnixNano(), 10) + str.String()
}

func BenchmarkFlashDB_Set(b *testing.B) {
	db := getTestDB()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.View(func(tx *Tx) error {
			err := tx.Set(newKey(i), randomValue())
			if err != nil {
				panic(err)
			}
			return nil
		})
	}
}

func BenchmarkFlashDB_Get(b *testing.B) {
	db := getTestDB()

	db.View(func(tx *Tx) error {
		for i := 0; i < b.N; i++ {
			tx.Set(newKey(i), randomValue())
		}
		return nil
	})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.View(func(tx *Tx) error {
			_, err := tx.Get(newKey(i))
			if err != nil && err != ErrInvalidKey {
				panic(err)
			}
			return nil
		})
	}
}

func BenchmarkFlashDB_HSet(b *testing.B) {
	db := getTestDB()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.View(func(tx *Tx) error {
			_, err := tx.HSet(newKey(i), randomValue(), randomValue())
			if err != nil {
				panic(err)
			}
			return nil
		})
	}
}

func BenchmarkFlashDB_HGet(b *testing.B) {
	db := getTestDB()

	db.View(func(tx *Tx) error {
		for i := 0; i < b.N; i++ {
			_, err := tx.HSet(newKey(i), newValue(i), randomValue())
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db.View(func(tx *Tx) error {
			val := tx.HGet(newKey(i), newValue(i))
			if val == "" {
				panic("empty value")
			}
			return nil
		})
	}
}
