# flashdb

FlashDB is a simple, in-memory, key/value store in pure Go.
It persists to disk, is ACID compliant, and uses locking for multiple
readers and a single writer. It supports redis like operations for
data structures like SET, SORTED SET, HASH and STRING. 

Features
========

- In-memory database for [fast reads and writes](#performance)
- Embeddable with a simple API
- Supports Redis like operatios for SET, SORTED SET, HASH and STRING
- [Durable append-only file](#append-only-file) format for persistence
- Option to evict old items with an [expiration](#data-expiration) TTL
- ACID semantics with locking [transactions](#transactions) that support rollbacks


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
    config:=&flashdb.Config{Path:"/tmp", EvictionInterval: 10}
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
config:=&flashdb.Config{Path:"", EvictionInterval: 10}
flashdb.New(config)
```

## Transactions
All reads and writes must be performed from inside a transaction. FlashDB can have one write transaction opened at a time, but can have many concurrent read transactions. Each transaction maintains a stable view of the database. In other words, once a transaction has begun, the data for that transaction cannot be changed by other transactions.

Transactions run in a function that exposes a `Tx` object, which represents the transaction state. While inside a transaction, all database operations should be performed using this object. You should never access the origin `DB` object while inside a transaction. Doing so may have side-effects, such as blocking your application.

When a transaction fails, it will roll back, and revert all changes that occurred to the database during that transaction. There's a single return value that you can use to close the transaction. For read/write transactions, returning an error this way will force the transaction to roll back. When a read/write transaction succeeds all changes are persisted to disk.

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

## Setting and getting key/values

To set a value you must open a read/write transaction:

```go
err := db.Update(func(tx *flashdb.Tx) error {
	_, _, err := tx.Set("mykey", "myvalue", nil)
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

Reference
=========

FlashDB is inspired by NutsDB and BuntDB.

