package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/arriqaaq/flashdb"
	"github.com/tidwall/redcon"
)

// ErrSyntaxIncorrect incorrect err
var ErrSyntaxIncorrect = errors.New("syntax err")
var okResult = redcon.SimpleString("OK")

func newWrongNumOfArgsError(cmd string) error {
	return fmt.Errorf("wrong number of arguments for '%s' command", cmd)
}

func set(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("set")
		return
	}

	key, value := args[0], args[1]
	if err := db.Update(func(tx *flashdb.Tx) error {
		return tx.Set(key, value)
	}); err == nil {
		res = okResult
	}

	return
}

func get(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("get")
		return
	}

	key := args[0]

	err = db.View(func(tx *flashdb.Tx) error {
		res, err = tx.Get(key)
		return nil
	})

	return
}

func delete(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("get")
		return
	}

	key := args[0]

	err = db.Update(func(tx *flashdb.Tx) error {
		return tx.Delete(key)
	})

	return
}

func expire(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = ErrSyntaxIncorrect
		return
	}
	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}

	if err := db.Update(func(tx *flashdb.Tx) error {
		return tx.Expire(args[0], int64(seconds))
	}); err == nil {
		res = okResult
	}

	return
}

func ttl(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("ttl")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		ttlVal := tx.TTL(args[0])
		res = strconv.FormatInt(int64(ttlVal), 10)
		return nil
	})

	return
}

func init() {
	addExecCommand("set", set)
	addExecCommand("get", get)
	addExecCommand("expire", expire)
	addExecCommand("ttl", ttl)
	addExecCommand("delete", delete)
}
