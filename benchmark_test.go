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
		err := db.Set(newKey(i), randomValue())
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFlashDB_Get(b *testing.B) {
	db := getTestDB()

	for i := 0; i < b.N; i++ {
		db.Set(newKey(i), randomValue())
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := db.Get(newKey(i))
		if err != nil && err != ErrInvalidKey {
			panic(err)
		}
	}
}

func BenchmarkFlashDB_HSet(b *testing.B) {
	db := getTestDB()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.HSet(newKey(i), randomValue(), randomValue())
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFlashDB_HGet(b *testing.B) {
	db := getTestDB()

	for i := 0; i < b.N; i++ {
		_, err := db.HSet(newKey(i), newValue(i), randomValue())
		if err != nil {
			panic(err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		val := db.HGet(newKey(i), newValue(i))
		if val == "" {
			panic("empty value")
		}
	}
}
