<p align="center">
<img
    src="img/architecture.png" alt="FlashDB">
</p>

# flashdb

FlashDB is a simple, in-memory, key/value store in pure Go.
It persists to disk, is ACID compliant, and uses locking for multiple
readers and a single writer. It supports redis like operations for
data structures like SET, SORTED SET, HASH and STRING. 

Features
========

- In-memory database for [fast reads and writes](#performance)
- Embeddable with a simple API
- Supports Redis like operations for SET, SORTED SET, HASH and STRING
- [Durable append-only file](#append-only-file) format for persistence
- Option to evict old items with an [expiration](#data-expiration) TTL
- ACID semantics with locking [transactions](#transactions) that support rollbacks


Architecture
=============

FlashDB is made of composable libraries that can be used independently and are easy to understand. The idea is to bridge the 
learning for anyone new on how to build a simple ACID database.


- [Set](https://github.com/arriqaaq/set)
- [ZSet](https://github.com/arriqaaq/zset)
- [String](https://github.com/arriqaaq/skiplist)
- [Hash](https://github.com/arriqaaq/hash)
- [Append Only Log](https://github.com/arriqaaq/aol)


Getting Started
===============

## Installing

To start using FlashDB, install Go and run `go get`:

```sh
$ go get -u github.com/arriqaaq/flashdb
```

This will retrieve the library.


## Opening a database

The primary object in FlashDB is a `DB`. To open or create your
database, use the `flashdb.New()` function:

```go
package main

import (
	"log"

	"github.com/arriqaaq/flashdb"
)

func main() {
    config := &flashdb.Config{Path:"/tmp", EvictionInterval: 10}
	db, err := flashdb.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	...
}
```

It's also possible to open a database that does not persist to disk by keeping the path in the config empty.

```go
config := &flashdb.Config{Path:"", EvictionInterval: 10}
flashdb.New(config)
```

## Transactions
All reads and writes must be performed from inside a transaction. FlashDB can have one write transaction opened at a time, but can have many concurrent read transactions. Each transaction maintains a stable view of the database. In other words, once a transaction has begun, the data for that transaction cannot be changed by other transactions.

When a transaction fails, it will roll back, and revert all changes that occurred to the database during that transaction. When a read/write transaction succeeds all changes are persisted to disk.

### Read-only Transactions
A read-only transaction should be used when you don't need to make changes to the data. The advantage of a read-only transaction is that there can be many running concurrently.

```go
err := db.View(func(tx *flashdb.Tx) error {
	...
	return nil
})
```

### Read/write Transactions
A read/write transaction is used when you need to make changes to your data. There can only be one read/write transaction running at a time. So make sure you close it as soon as you are done with it.

```go
err := db.Update(func(tx *flashdb.Tx) error {
	...
	return nil
})
```

### Setting and getting key/values

To set a value you must open a read/write transaction:

```go
err := db.Update(func(tx *flashdb.Tx) error {
	_, _, err := tx.Set("mykey", "myvalue")
	return err
})
```


To get the value:

```go
err := db.View(func(tx *flashdb.Tx) error {
	val, err := tx.Get("mykey")
	if err != nil{
		return err
	}
	fmt.Printf("value is %s\n", val)
	return nil
})
```

Commands
========
| String | Hash    | Set         | ZSet           |
|--------|---------|-------------|----------------|
| SET    | HSET    | SADD        | ZADD           |
| GET    | HGET    | SISMEMBER   | ZSCORE         |
| DELETE | HGETALL | SRANDMEMBER | ZCARD          |
| EXPIRE | HDEL    | SREM        | ZRANK          |
| TTL    | HEXISTS | SMOVE       | ZREVRANK       |
|        | HLEN    | SCARD       | ZRANGE         |
|        | HKEYS   | SMEMBERS    | ZREVRANGE      |
|        | HVALS   | SUNION      | ZREM           |
|        | HCLEAR  | SDIFF       | ZGETBYRANK     |
|        |         | SCLEAR      | ZREVGETBYRANK  |
|        |         |             | ZSCORERANGE    |
|        |         |             | ZREVSCORERANGE |
|        |         |             | ZCLEAR         |



TODO
====

FlashDB is a toy db and in early stages of development. A couple of to-do tasks listed:

- Add transaction to eviction to avoid locks
- Add more comprehensive unit test cases
- Add more bench tests
- Add explicit logging messages


References
==========

FlashDB is inspired by NutsDB and BuntDB.


## Contact
Farhan Khan [@arriqaaq](http://twitter.com/arriqaaq)

## License
FlashDB source code is available under the MIT [License](/LICENSE)
