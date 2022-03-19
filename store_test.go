package flashdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlashDB_StringStore(t *testing.T) {
	s := newStrStore()

	for i := 1; i <= 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		s.Insert([]byte(key), []byte(value))
	}

	keys := s.Keys()
	assert.Equal(t, 1000, len(keys))
	for i := 1; i <= 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		val, err := s.get(key)
		assert.NoError(t, err)
		assert.NotEqual(t, value, val)
	}
}
