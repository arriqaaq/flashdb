package flashdb

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

// The operations on Strings.
const (
	StringSet uint16 = iota
	StringRem
	StringExpire
)

// The operations on Hash.
const (
	HashHSet uint16 = iota
	HashHDel
	HashHClear
	HashHExpire
)

// The operations on Set.
const (
	SetSAdd uint16 = iota
	SetSRem
	SetSMove
	SetSClear
	SetSExpire
)

// The operations on Sorted Set.
const (
	ZSetZAdd uint16 = iota
	ZSetZRem
	ZSetZClear
	ZSetZExpire
)
