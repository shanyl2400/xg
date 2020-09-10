package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"xg/log"
)

const (
	// DBTransactionTimeout transaction timeout
	DBTransactionTimeout = time.Second * 8
)

func getDBTransactionTimeout() time.Duration {
	return DBTransactionTimeout
}

// GetTrans begin a transaction
func GetTrans(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	log.Info.Println("begin transaction")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, getDBTransactionTimeout())
	defer cancel()

	tx := Get().BeginTx(ctxWithTimeout, &sql.TxOptions{})

	funcDone := make(chan error, 0)
	var err error
	go func() {
		defer func() {
			if err1 := recover(); err1 != nil {
				log.Error.Printf("recover error: %v with transaction panic", err1)
				funcDone <- fmt.Errorf("transaction panic: %+v", err1)
			}
		}()

		// call func
		funcDone <- fn(ctxWithTimeout, tx)
	}()

	select {
	case err = <-funcDone:
		log.Trace.Println("transaction fn done")
	case <-ctxWithTimeout.Done():
		// context deadline exceeded
		err = ctxWithTimeout.Err()
		log.Error.Println("transaction context deadline exceeded, ", err)
	}

	if err != nil {
		err1 := tx.RollbackUnlessCommitted().Error
		if err1 != nil {
			log.Error.Println("rollback transaction failed, ", err, ", err1:", err1)
		} else {
			log.Trace.Println("rollback transaction success")
		}
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error.Println("commit transaction failed, ", err)
		return err
	}
	log.Trace.Println("commit transaction success")
	return nil
}

type transactionResult struct {
	Result interface{}
	Error  error
}

// GetTransResult begin a transaction, get result of callback
func GetTransResult(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	log.Info.Println("begin transaction")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, getDBTransactionTimeout())
	defer cancel()


	tx := Get().BeginTx(ctxWithTimeout, &sql.TxOptions{})
	var err error

	funcDone := make(chan *transactionResult, 0)
	go func() {
		defer func() {
			if err1 := recover(); err1 != nil {
				log.Error.Printf("recover error: %v with transaction panic", err1)
				funcDone <- &transactionResult{Error: fmt.Errorf("transaction panic: %+v", err1)}
			}
		}()

		// call func
		result, err := fn(ctxWithTimeout, tx)
		funcDone <- &transactionResult{Result: result, Error: err}
	}()

	var funcResult *transactionResult
	select {
	case funcResult = <-funcDone:
		log.Trace.Println("transaction fn done")
	case <-ctxWithTimeout.Done():
		// context deadline exceeded
		funcResult = &transactionResult{Error: ctxWithTimeout.Err()}
		err = ctxWithTimeout.Err()
		log.Error.Println("transaction context deadline exceeded, ", err)
	}

	if funcResult.Error != nil {
		log.Error.Println("transaction failed, ", funcResult.Error)
		err1 := tx.RollbackUnlessCommitted().Error
		if err1 != nil {
			log.Error.Println("rollback transaction failed, ", err, ", err1:", err1)
		} else {
			log.Trace.Println("rollback transaction success")
		}
		return nil, funcResult.Error
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error.Println("commit transaction failed, ", err)
		return nil, err
	}

	log.Trace.Println("commit transaction success")

	return funcResult.Result, nil
}
