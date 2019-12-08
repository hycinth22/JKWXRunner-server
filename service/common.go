package service

import (
	"sync"

	"github.com/inkedawn/JKWXRunner-server/database"
)

// comprised of every service interface
type ICommonService interface {
	// Begin a transaction and return it.
	// The return value is only for r/w operation in raw transaction. (for example, external library or old style service)
	// Don't use database.TX.Commit/Rollback. Use ICommonService.Begin/Rollback instead.
	Begin() (tx database.TX)
	// Commit the transaction on this
	Commit()
	// Rollback the transaction on this
	Rollback()
	// underlying database, without any transaction.
	GetDB() (txOrDb *database.DB)
}

type commonService struct {
	db *database.DB
	tx *database.DB

	cntTransaction uint // support nested transaction feature
	m              sync.Locker
}

func (c *commonService) Begin() database.TX {
	c.m.Lock()
	defer c.m.Unlock()
	if c.cntTransaction == 0 {
		c.tx = c.db.Begin()
	}
	c.cntTransaction++
	return c.tx
}

func (c *commonService) Commit() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.cntTransaction == 0 {
		return // no transaction
	}
	c.cntTransaction--
	if c.cntTransaction > 0 {
		return // commit nested transaction
	}
	// final commit
	c.tx.Commit()
	c.tx = nil
	return
}

// notice: it's not a true nested transaction. nested transaction rollback operation was dropped.
// work improperly if nested transaction rollback but parent transaction commit
func (c *commonService) Rollback() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.cntTransaction == 0 {
		return // no transaction
	}
	c.cntTransaction--
	if c.cntTransaction > 0 {
		return // rollback nested transaction.
	}
	// final rollback
	c.tx.Rollback()
	c.tx = nil
	return
}

func (c *commonService) GetDB() *database.DB {
	if c.tx != nil {
		return c.tx
	}
	return c.db
}

func NewCommonService() ICommonService {
	return NewCommonServiceOn(database.GetDB())
}

func NewCommonServiceOn(db *database.DB) ICommonService {
	return &commonService{db: db, m: &sync.Mutex{}}
}
