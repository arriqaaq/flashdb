package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arriqaaq/flashdb"
	"github.com/tidwall/redcon"
)

// float64ToStr Convert type float64 to string
func float64ToStr(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

// strToFloat64 convert type string to float64
func strToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func zAdd(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 {
		err = newWrongNumOfArgsError("zadd")
		return
	}
	score, err := strToFloat64(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.Update(func(tx *flashdb.Tx) error {
		return tx.ZAdd(args[0], score, args[2])
	})

	if err == nil {
		res = okResult
	}
	return
}

func zScore(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("zscore")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		ok, score := tx.ZScore(args[0], args[1])
		if ok {
			res = float64ToStr(score)
		}
		return nil
	})

	return
}

func zCard(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("zcard")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		card := tx.ZCard(args[0])
		res = redcon.SimpleInt(card)
		return nil
	})

	return
}

func zRank(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		rank := tx.ZRank(args[0], args[1])
		res = redcon.SimpleInt(rank)
		return nil
	})

	return
}

func zRevRank(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("zrevrank")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		rank := tx.ZRevRank(args[0], args[1])
		res = redcon.SimpleInt(rank)
		return nil
	})

	return
}

func zRange(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 && len(args) != 4 {
		err = newWrongNumOfArgsError("zrange")
		return
	}
	return zRawRange(db, args, false)
}

func zRevRange(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 && len(args) != 4 {
		err = newWrongNumOfArgsError("zrevrange")
		return
	}
	return zRawRange(db, args, true)
}

// for zRange and zRevRange
func zRawRange(db *flashdb.FlashDB, args []string, rev bool) (res interface{}, err error) {
	withScores := false
	if len(args) == 4 {
		if strings.ToLower(args[3]) == "withscores" {
			withScores = true
			args = args[:3]
		} else {
			err = ErrSyntaxIncorrect
			return
		}
	}
	start, err := strconv.Atoi(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}

	var val []interface{}
	err = db.View(func(tx *flashdb.Tx) error {
		if rev {
			if withScores {
				val = tx.ZRevRangeWithScores(args[0], start, end)
			} else {
				val = tx.ZRevRange(args[0], start, end)
			}
		} else {
			if withScores {
				val = tx.ZRangeWithScores(args[0], start, end)
			} else {
				val = tx.ZRange(args[0], start, end)
			}
		}
		return nil
	})

	results := make([]string, len(val))
	for i, v := range val {
		results[i] = fmt.Sprintf("%v", v)
	}
	res = results
	return
}

func zRem(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = ErrSyntaxIncorrect
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		var ok bool
		if ok, err = tx.ZRem(args[0], args[1]); err == nil {
			if ok {
				res = redcon.SimpleInt(1)
			} else {
				res = redcon.SimpleInt(0)
			}
		}
		return nil
	})

	return
}

func zGetByRank(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("zgetbyrank")
		return
	}
	return zRawGetByRank(db, args, false)
}

func zRevGetByRank(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 2 {
		err = newWrongNumOfArgsError("zrevgetbyrank")
		return
	}
	return zRawGetByRank(db, args, true)
}

// for zGetByRank and zRevGetByRank
func zRawGetByRank(db *flashdb.FlashDB, args []string, rev bool) (res interface{}, err error) {
	rank, err := strconv.Atoi(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}

	var val []interface{}

	err = db.View(func(tx *flashdb.Tx) error {
		if rev {
			val = tx.ZRevGetByRank(args[0], rank)
		} else {
			val = tx.ZGetByRank(args[0], rank)
		}
		return nil
	})

	results := make([]string, len(val))
	for i, v := range val {
		results[i] = fmt.Sprintf("%v", v)
	}
	res = results
	return
}

func zScoreRange(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 {
		err = newWrongNumOfArgsError("zscorerange")
		return
	}
	return zRawScoreRange(db, args, false)
}

func zSRevScoreRange(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 3 {
		err = newWrongNumOfArgsError("zsrevscorerange")
		return
	}
	return zRawScoreRange(db, args, true)
}

// for zScoreRange and zSRevScoreRange
func zRawScoreRange(db *flashdb.FlashDB, args []string, rev bool) (res interface{}, err error) {
	param1, err := strToFloat64(args[1])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}
	param2, err := strToFloat64(args[2])
	if err != nil {
		err = ErrSyntaxIncorrect
		return
	}
	var val []interface{}

	err = db.View(func(tx *flashdb.Tx) error {
		if rev {
			val = tx.ZRevScoreRange(args[0], param1, param2)
		} else {
			val = tx.ZScoreRange(args[0], param1, param2)
		}
		return nil
	})

	results := make([]string, len(val))
	for i, v := range val {
		results[i] = fmt.Sprintf("%v", v)
	}
	res = results
	return
}

func zclear(db *flashdb.FlashDB, args []string) (res interface{}, err error) {
	if len(args) != 1 {
		err = newWrongNumOfArgsError("zclear")
		return
	}

	err = db.View(func(tx *flashdb.Tx) error {
		if err := tx.ZClear(args[0]); err == nil {
			res = redcon.SimpleInt(1)
		} else {
			res = redcon.SimpleInt(0)
		}
		return nil
	})

	return
}

func init() {
	addExecCommand("zadd", zAdd)
	addExecCommand("zscore", zScore)
	addExecCommand("zcard", zCard)
	addExecCommand("zrank", zRank)
	addExecCommand("zrevrank", zRevRank)
	addExecCommand("zrange", zRange)
	addExecCommand("zrevrange", zRevRange)
	addExecCommand("zrem", zRem)
	addExecCommand("zgetbyrank", zGetByRank)
	addExecCommand("zrevgetbyrank", zRevGetByRank)
	addExecCommand("zscorerange", zScoreRange)
	addExecCommand("zrevscorerange", zSRevScoreRange)
	addExecCommand("zclear", zclear)
}
