package repository

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

type Model interface {
	TableName() string
}

type Models []Model

type Cuder interface {
	Transact(func(*gorm.DB) error) error
	Insert(Model) error
	Update(Model, ColumnsAndValues) error
	UpdateWithWhere(Model, ColumnsAndValues, FilterCondition) (int64, error)
	Delete(Model) error
	DeleteWithWhere(Model, FilterCondition) error
}

type cuder struct {
	db *gorm.DB
}

type RawFilter struct {
	RawSQL       string
	RawSQLParams []string
}

func (rf *RawFilter) Validate() bool {
	c := strings.Count(rf.RawSQL, "?")
	return c == len(rf.RawSQLParams)
}

type FilterCondition struct {
	FilterColumnsAndValues map[string]interface{}
	RawFilter              *RawFilter
}

type ColumnsAndValues map[string]interface{}

func NewCuder(db *gorm.DB) Cuder {
	return cuder{db}
}

func (c cuder) Insert(m Model) error {
	return c.db.Create(m).Error
}

func (c cuder) Update(m Model, cvs ColumnsAndValues) error {
	return c.db.Model(m).Updates(cvs).Error
}

func (c cuder) UpdateWithWhere(m Model, cvs ColumnsAndValues, fc FilterCondition) (int64, error) {
	db := c.db.Model(m)

	if fc.FilterColumnsAndValues != nil {
		db = db.Where(fc.FilterColumnsAndValues)
	} else if fc.RawFilter != nil {
		ok := fc.RawFilter.Validate()
		if !ok {
			return 0, fmt.Errorf("Invalid RawFilter passed")
		}
		db.Where(fc.RawFilter.RawSQLParams, fc.RawFilter.RawSQLParams)
	}
	res := db.Updates(cvs)
	return res.RowsAffected, res.Error
}

func (c cuder) Delete(m Model) error {
	return c.db.Delete(m).Error
}

func (c cuder) DeleteWithWhere(m Model, fc FilterCondition) error {
	db := c.db.Model(m)

	if fc.FilterColumnsAndValues != nil {
		db = db.Where(fc.FilterColumnsAndValues)
	} else if fc.RawFilter != nil {
		ok := fc.RawFilter.Validate()
		if !ok {
			return fmt.Errorf("Invalid RawFilter passed")
		}
		db.Where(fc.RawFilter.RawSQLParams, fc.RawFilter.RawSQLParams)
	}
	return db.Delete(m).Error
}

func (c cuder) Transact(txFunc func(*gorm.DB) error) (err error) {
	tx := c.db.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			err = fmt.Errorf("Panic: %v", p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = txFunc(tx)

	return err
}
