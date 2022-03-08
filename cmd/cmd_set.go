package cmd

import (
	"strconv"

	"github.com/arriqaaq/flashdb"
	"github.com/tidwall/redcon"
)

func sAdd(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) <= 1 {
		err = newWrongNumOfArgsError("sadd")
		return
	}

	var members []string
	members = append(members, args[1:]...)

	err = db.View(func(tx *flashdb.Tx) error {
		err = tx.SAdd(args[0], members...)
		return err
	})

	return
}

func sIsMember(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("sismember")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		if ok := tx.SIsMember(args[0], args[1]); ok {
			res = redcon.SimpleInt(1)
		} else {
			res = redcon.SimpleInt(0)
		}
		return nil
	})

	return
}

func sRandMember(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("srandmember")
		return
	}
	count, err := strconv.Atoi(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.SRandMember(args[0], count)
		return nil
	})

	return
}

func sRem(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) <= 1 {
		err = ErrSyntaxIncorrect
		return
	}
	var members []string
	members = append(members, args[1:]...)

	err = db.View(func(tx *flashdb.Tx) error {
		var count int
		if count, err = tx.SRem(args[0], members...); err == nil {
			res = redcon.SimpleInt(count)
		}
		return err
	})

	return
}

func sMove(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 {
		err = newWrongNumOfArgsError("smove")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		return tx.SMove(args[0], args[1], args[2])
	})

	if err == nil {
		res = okResult
	}

	return
}

func sCard(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("scard")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		card := tx.SCard(args[0])
		res = redcon.SimpleInt(card)
		return nil
	})

	return
}

func sMembers(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("smembers")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.SMembers(args[0])
		return nil
	})

	return
}

func sUnion(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) <= 0 {
		err = newWrongNumOfArgsError("sunion")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.SUnion(args...)
		return nil
	})

	return
}

func sDiff(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) <= 0 {
		err = newWrongNumOfArgsError("sdiff")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		res = tx.SDiff(args...)
		return nil
	})

	return
}

func sclear(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("sclear")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		if err := tx.SClear(args[0]); err == nil {
			res = redcon.SimpleInt(1)
		} else {
			res = redcon.SimpleInt(0)
		}
		return nil
	})

	return
}

func init() {
	addExecCommand("sadd", sAdd)
	addExecCommand("sismember", sIsMember)
	addExecCommand("srandmember", sRandMember)
	addExecCommand("srem", sRem)
	addExecCommand("smove", sMove)
	addExecCommand("scard", sCard)
	addExecCommand("smembers", sMembers)
	addExecCommand("sunion", sUnion)
	addExecCommand("sdiff", sDiff)
	addExecCommand("sclear", sclear)
}
