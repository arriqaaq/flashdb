package flashdb

import (
	"github.com/arriqaaq/aol"
)

// Tx represents a transaction on the database. This transaction can either be
// read-only or read/write. Read-only transactions can be used for retrieving
// values for keys and iterating through keys and values. Read/write
// transactions can set and delete keys.
//
// All transactions must be committed or rolled-back when done.
type Tx struct {
	db       *FlashDB        // the underlying database.
	writable bool            // when false mutable operations fail.
	wc       *txWriteContext // context for writable transactions.
}

func (tx *Tx) addRecord(r *record) {
	if tx.db.persist {
		tx.wc.rollbackItems = append(tx.wc.rollbackItems, r)
		tx.wc.commitItems = append(tx.wc.commitItems, r)
	}
}

type txWriteContext struct {
	rollbackItems []*record // details for rolling back tx.
	commitItems   []*record // details for committing tx.
}

// lock locks the database based on the transaction type.
func (tx *Tx) lock() {
	if tx.writable {
		tx.db.mu.Lock()
	} else {
		tx.db.mu.RLock()
	}
}

// unlock unlocks the database based on the transaction type.
func (tx *Tx) unlock() {
	if tx.writable {
		tx.db.mu.Unlock()
	} else {
		tx.db.mu.RUnlock()
	}
}

// managed calls a block of code that is fully contained in a transaction.
// This method is intended to be wrapped by Update and View
func (db *FlashDB) managed(writable bool, fn func(tx *Tx) error) (err error) {
	var tx *Tx
	tx, err = db.Begin(writable)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			// The caller returned an error. We must rollback.
			_ = tx.Rollback()
			return
		}
		if writable {
			// Everything went well. Lets Commit()
			err = tx.Commit()
		} else {
			// read-only transaction can only roll back.
			err = tx.Rollback()
		}
	}()
	err = fn(tx)
	return
}

// View executes a function within a managed read-only transaction.
// When a non-nil error is returned from the function that error will be return
// to the caller of View().
func (db *FlashDB) View(fn func(tx *Tx) error) error {
	return db.managed(false, fn)
}

// Update executes a function within a managed read/write transaction.
// The transaction has been committed when no error is returned.
// In the event that an error is returned, the transaction will be rolled back.
// When a non-nil error is returned from the function, the transaction will be
// rolled back and the that error will be return to the caller of Update().
func (db *FlashDB) Update(fn func(tx *Tx) error) error {
	return db.managed(true, fn)
}

// Rollback closes the transaction and reverts all mutable operations that
// were performed on the transaction such as Set() and Delete().
//
// Read-only transactions can only be rolled back, not committed.
func (tx *Tx) Rollback() error {
	if tx.db == nil {
		return ErrTxClosed
	}
	// The rollback func does the heavy lifting.
	if tx.writable {
		tx.rollbackInner()
	}
	// unlock the database for more transactions.
	tx.unlock()
	// Clear the db field to disable this transaction from future use.
	tx.db = nil
	return nil
}

// rollbackInner handles the underlying rollback logic.
// Intended to be called from Commit() and Rollback().
func (tx *Tx) rollbackInner() {
	for _, item := range tx.wc.rollbackItems {
		tx.db.rollbackFromDatabase(item)
	}
}

// Commit writes all changes to disk.
// An error is returned when a write error occurs, or when a Commit() is called
// from a read-only transaction.
func (tx *Tx) Commit() error {
	if tx.db == nil {
		return ErrTxClosed
	} else if !tx.writable {
		return ErrTxNotWritable
	}
	var err error
	if tx.db.persist && (len(tx.wc.commitItems) > 0) {
		batch := new(aol.Batch)
		// Each committed record is written to disk
		for _, r := range tx.wc.commitItems {
			rec, err := r.encode()
			if err != nil {
				return err
			}
			batch.Write(rec)
		}
		// If this operation fails then the write did failed and we must
		// rollback.
		err = tx.db.log.WriteBatch(batch)
		if err != nil {
			tx.rollbackInner()
		}
	}
	// Unlock the database and allow for another writable transaction.
	tx.unlock()
	// Clear the db field to disable this transaction from future use.
	tx.db = nil
	return err
}

// Begin opens a new transaction.
// Multiple read-only transactions can be opened at the same time but there can
// only be one read/write transaction at a time. Attempting to open a read/write
// transactions while another one is in progress will result in blocking until
// the current read/write transaction is completed.
//
// All transactions must be closed by calling Commit() or Rollback() when done.
func (db *FlashDB) Begin(writable bool) (*Tx, error) {
	tx := &Tx{
		db:       db,
		writable: writable,
	}
	tx.lock()
	if db.closed {
		tx.unlock()
		return nil, ErrDatabaseClosed
	}
	if writable {
		tx.wc = &txWriteContext{}
		tx.wc.rollbackItems = make([]*record, 0)
		if db.persist {
			tx.wc.commitItems = make([]*record, 0)
		}
	}
	return tx, nil
}

// rollbackFromDatabase removes and item from the database and expiry list.
func (db *FlashDB) rollbackFromDatabase(r *record) {
	key := string(r.meta.key)
	switch r.getType() {
	case StringRecord:
		db.strStore.Delete(key)
		db.exps.HDel(String, key)
	case HashRecord:
		db.hashStore.HClear(key)
		db.exps.HDel(Hash, key)
	case SetRecord:
		db.setStore.SClear(key)
		db.exps.HDel(Set, key)
	case ZSetRecord:
		db.zsetStore.ZClear(key)
		db.exps.HDel(ZSet, key)
	}
}
