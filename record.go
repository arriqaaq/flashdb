package flashdb

import (
	"encoding/binary"
	"errors"
	"time"
)

/*
	Record is based on bitcask entry model design in nutsdb
*/

var (
	ErrInvalidEntry = errors.New("invalid entry")
	ErrInvalidCrc   = errors.New("invalid crc")
)

const (
	// keySize, memberSize, valueSize is uint32 type，4 bytes each.
	// timestamp 8 bytes, state 2 bytes.
	// 4 + 4 + 4 + 8 + 2 = 30
	entryHeaderSize = 22
)

type (
	record struct {
		meta      *meta
		state     uint16 // state represents two fields, high 8 bits is the data type, low 8 bits is operation mark.
		timestamp uint64 // Timestamp is the time when entry was written.
	}

	// Meta meta info.
	meta struct {
		key        []byte
		member     []byte
		value      []byte
		keySize    uint32
		memberSize uint32
		valueSize  uint32
	}
)

func newInternal(key, member, value []byte, state uint16, timestamp uint64) *record {
	return &record{
		state: state, timestamp: timestamp,
		meta: &meta{
			key:        key,
			member:     member,
			value:      value,
			keySize:    uint32(len(key)),
			memberSize: uint32(len(member)),
			valueSize:  uint32(len(value)),
		},
	}
}

func newRecord(key, member []byte, t, mark uint16) *record {
	var state uint16 = 0
	// set type and mark.
	state = state | (t << 8)
	state = state | mark
	return newInternal(key, member, nil, state, uint64(time.Now().UnixNano()))
}

func newRecordWithValue(key, member, value []byte, t, mark uint16) *record {
	var state uint16 = 0
	// set type and mark.
	state = state | (t << 8)
	state = state | mark
	return newInternal(key, member, value, state, uint64(time.Now().UnixNano()))
}

func newRecordWithExpire(key, member []byte, deadline int64, t, mark uint16) *record {
	var state uint16 = 0
	// set type and mark.
	state = state | (t << 8)
	state = state | mark

	return newInternal(key, member, nil, state, uint64(deadline))
}

func (e *record) size() uint32 {
	return entryHeaderSize + e.meta.keySize + e.meta.memberSize + e.meta.valueSize
}

// Encode returns the slice after the entry be encoded.
//
//  the entry stored format:
//  |----------------------------------------------------------------------------------------------------------------|
//  |  ks   | ms      | vs     | state  | timestamp  | key    | member | value  |
//  |----------------------------------------------------------------------------------------------------------------|
//  | uint32| uint32  | uint32 | uint16 | uint64     | []byte | []byte | []byte |
//  |----------------------------------------------------------------------------------------------------------------|
//

func (e *record) encode() ([]byte, error) {
	if e == nil || e.meta.keySize == 0 {
		return nil, ErrInvalidEntry
	}

	ks, ms := e.meta.keySize, e.meta.memberSize
	vs := e.meta.valueSize
	buf := make([]byte, e.size())

	binary.BigEndian.PutUint32(buf[0:4], ks)
	binary.BigEndian.PutUint32(buf[4:8], ms)
	binary.BigEndian.PutUint32(buf[8:12], vs)
	binary.BigEndian.PutUint16(buf[12:14], e.state)
	binary.BigEndian.PutUint64(buf[14:22], e.timestamp)
	copy(buf[entryHeaderSize:entryHeaderSize+ks], e.meta.key)
	copy(buf[entryHeaderSize+ks:(entryHeaderSize+ks+ms)], e.meta.member)
	if vs > 0 {
		copy(buf[(entryHeaderSize+ks+ms):(entryHeaderSize+ks+ms+vs)], e.meta.value)
	}

	return buf, nil
}

func decode(buf []byte) (*record, error) {
	ks := binary.BigEndian.Uint32(buf[0:4])
	ms := binary.BigEndian.Uint32(buf[4:8])
	vs := binary.BigEndian.Uint32(buf[8:12])
	state := binary.BigEndian.Uint16(buf[12:14])
	timestamp := binary.BigEndian.Uint64(buf[14:22])

	return &record{
		meta: &meta{
			keySize:    ks,
			memberSize: ms,
			valueSize:  vs,
			key:        buf[entryHeaderSize : entryHeaderSize+ks],
			member:     buf[entryHeaderSize+ks : (entryHeaderSize + ks + ms)],
			value:      buf[(entryHeaderSize + ks + ms):(entryHeaderSize + ks + ms + vs)],
		},
		state:     state,
		timestamp: timestamp,
	}, nil
}

func (e *record) getType() uint16 {
	return e.state >> 8
}

func (e *record) getMark() uint16 {
	return e.state & (2<<7 - 1)
}
