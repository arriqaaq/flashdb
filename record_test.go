package flashdb

import (
	"fmt"
	"os"
	"testing"

	"github.com/arriqaaq/aol"
	"github.com/stretchr/testify/assert"
)

func makeRecords(n int) []*record {
	rec := make([]*record, 0, n)
	for i := 1; i <= n; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		member := fmt.Sprintf("member_%d", i)
		rec = append(rec, newRecordWithValue([]byte(key), []byte(member), []byte(value), ZSetRecord, ZSetZAdd))
	}
	return rec
}

func TestFlashDB_AOL(t *testing.T) {
	logPath := "tmp/"
	l, err := aol.Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	defer os.RemoveAll(logPath)

	recs := makeRecords(100)
	for _, r := range recs {
		data, err := r.encode()
		assert.NoError(t, err)
		l.Write(data)
	}

	var lastRecord *record

	segs := l.Segments()
	for i := 1; i <= segs; i++ {
		j := 0
		for {
			data, err := l.Read(uint64(i), uint64(j))
			if err != nil {
				if err == aol.ErrEOF {
					break
				}
				t.Fatalf("expected %v, got %v", nil, err)
			}
			res, err := decode(data)
			assert.NoError(t, err)
			lastRecord = res
			j++
		}
	}

	assert.Equal(t, "key_100", string(lastRecord.meta.key))
	assert.Equal(t, "member_100", string(lastRecord.meta.member))
	assert.Equal(t, "value_100", string(lastRecord.meta.value))
}
