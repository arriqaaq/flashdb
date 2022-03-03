package flashdb

// DataType Define the data structure type.
type DataType = string

const (
	String DataType = "String"
	Hash   DataType = "Hash"
	Set    DataType = "Set"
	ZSet   DataType = "ZSet"
)

const (
	StringRecord uint16 = iota
	HashRecord
	SetRecord
	ZSetRecord
)

// The operations of a String Type, will be a part of Entry, the same for the other four types.
const (
	StringSet uint16 = iota
	StringRem
	StringExpire
)

// The operations of Hash.
const (
	HashHSet uint16 = iota
	HashHDel
	HashHClear
	HashHExpire
)

// The operations of Set.
const (
	SetSAdd uint16 = iota
	SetSRem
	SetSMove
	SetSClear
	SetSExpire
)

// The operations of Sorted Set.
const (
	ZSetZAdd uint16 = iota
	ZSetZRem
	ZSetZClear
	ZSetZExpire
)
