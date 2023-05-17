package dal

import (
	"flaver/globals"
	"flaver/lib/dal/database"
	"flaver/lib/dal/database/dal"

	"gorm.io/gorm"
)

type TransactionTask func(tx *gorm.DB) (interface{}, error)
type TransactionContext func(TransactionTask) (interface{}, error)

type DALWrapper struct {
	Dal                dal.IDal
	TransactionContext TransactionContext
}

func (this *DALWrapper) GetDal() dal.IDal {
	return this.Dal
}

func (this *DALWrapper) ExecTransaction(transaction TransactionTask) (interface{}, error) {

	return this.TransactionContext(transaction)
}

func NewDal() *DALWrapper {
	dal := dal.NewDal()
	return &DALWrapper{
		Dal:                dal,
		TransactionContext: NewTransactionContext(database.GetClientDB(), dal),
	}
}

func NewTransactionContext(db *gorm.DB, getSettables ...database.GormGetSettable) TransactionContext {
	return func(transaction TransactionTask) (result interface{}, err error) {
		// Initiate transaction
		tx := db.Begin()
		err = tx.Error
		if err != nil {
			return nil, err
		}
		defer func() {
			if r := recover(); r != nil {
				globals.GetLogger().Errorf("[NewTransactionContext] %+v", r)
				tx.Rollback()
			} else if err != nil {
				globals.GetLogger().Warnf("[NewTransactionContext] %+v", err)
				tx.Rollback()
			} else if err = tx.Commit().Error; err != nil {
				globals.GetLogger().Warnf("[NewTransactionContext] %+v", err)
			}
		}()
		// Store original DB connection for each getSettable and then set the tx as their DB connection
		originalDBs := make(map[int]*gorm.DB)
		for idx, getSettable := range getSettables {
			originalDBs[idx] = getSettable.GetConn()
			getSettable.SetConn(tx)
		}
		// Restore original DB connection for each getSettable while transaction has finished
		defer func() {
			for idx, db := range originalDBs {
				getSettables[idx].SetConn(db)
			}
		}()
		// Start transaction
		result, err = transaction(tx)
		return result, err
	}
}
