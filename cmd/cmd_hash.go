package cmd

import (
	"github.com/arriqaaq/flashdb"
	"github.com/tidwall/redcon"
)

func hSet(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 {
		err = newWrongNumOfArgsError("hset")
		return
	}

	err = db.Update(func(tx *flashdb.Tx) error {
		var count int
		if count, err = tx.HSet(args[0], args[1], args[2]); err == nil {
			res = redcon.SimpleInt(count)
		}
		return err
	})

	return
}

func hGet(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.HGet(args[0], args[1])
		return nil
	})

	return
}

func hGetAll(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("hgetall")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.HGetAll(args[0])
		return nil
	})

	return
}

func hDel(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) <= 1 {
		err = newWrongNumOfArgsError("hdel")
		return
	}

	var fields []string
	for _, f := range args[1:] {
		fields = append(fields, f)
	}

	err = db.Update(func(tx *flashdb.Tx) error {
		var count int
		if count, err = tx.HDel(args[0], fields...); err == nil {
			res = redcon.SimpleInt(count)
		}
		return err
	})

	return
}

func hExists(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("hexists")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		exists := tx.HExists(args[0], args[1])
		if exists {
			res = redcon.SimpleInt(1)
		} else {
			res = redcon.SimpleInt(0)
		}
		return nil
	})

	return
}

func hLen(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("hlen")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		count := tx.HLen(args[0])
		res = redcon.SimpleInt(count)
		return nil
	})

	return
}

func hKeys(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.HKeys(args[0])
		return nil
	})

	return
}

func hVals(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("hvals")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.HVals(args[0])
		return nil
	})

	return
}

func hClear(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.Update(func(tx *flashdb.Tx) error {
		res = tx.HClear(args[0])
		return nil
	})

	return
}

func init() {
	addExecCommand("hset", hSet)
	addExecCommand("hget", hGet)
	addExecCommand("hgetall", hGetAll)
	addExecCommand("hdel", hDel)
	addExecCommand("hexists", hExists)
	addExecCommand("hlen", hLen)
	addExecCommand("hkeys", hKeys)
	addExecCommand("hvals", hVals)
	addExecCommand("hclear", hClear)
}
