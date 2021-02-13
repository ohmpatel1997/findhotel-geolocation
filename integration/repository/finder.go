package repository

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Finder interface {
	FindByPKS(Model, ...interface{}) (Model, bool, error)
	FindByPKSStrict(Model, ...interface{}) (Model, error)
	FindManaged(Model) (Model, bool, error)
	FindManyManaged(interface{}, map[string]interface{}) (interface{}, error)
	FindRaw(interface{}, string, ...interface{}) (interface{}, bool, error)
	FindManyRaw(interface{}, string, ...interface{}) (interface{}, bool, error)
}

type finder struct {
	db *gorm.DB
}

func NewFinder(db *gorm.DB) Finder {
	return finder{db: db}
}

func (f finder) FindByPKS(m Model, pks ...interface{}) (Model, bool, error) {
	shouldFind := false

	for _, pk := range pks {

		// Check if UUID string
		s, ok := pk.(string)
		if ok {
			uID, err := uuid.Parse(s)
			if err != nil {
				//We have a non uuid pk
				if len(s) > 0 {
					shouldFind = true
					break
				}
			}

			if uID != uuid.Nil {
				shouldFind = true
				break
			}
		} else {
			//We have some other type, just check pk
			shouldFind = true
			break
		}
	}

	if !shouldFind {
		return m, false, nil
	}

	err := f.db.First(m, pks).Error
	if gorm.IsRecordNotFoundError(err) {
		return m, false, nil
	}

	if err != nil {
		return m, false, err
	}

	return m, true, nil
}

func (f finder) FindByPKSStrict(m Model, pks ...interface{}) (Model, error) {
	err := f.db.First(m, pks).Error

	return m, err
}

func (f finder) FindManaged(m Model) (Model, bool, error) {
	err := f.db.Where(m).First(m).Error

	if gorm.IsRecordNotFoundError(err) {
		// Let caller decide what to do
		return m, false, nil
	}

	if err != nil {
		return m, false, err
	}

	return m, true, nil
}

func (f finder) FindManyManaged(ms interface{}, cvs map[string]interface{}) (interface{}, error) {
	err := f.db.Where(cvs).Find(ms).Error

	return ms, err
}

func (f finder) FindRaw(m interface{}, querySQL string, params ...interface{}) (interface{}, bool, error) {
	err := f.db.Raw(querySQL, params...).Scan(m).Error

	if gorm.IsRecordNotFoundError(err) {
		// Let caller decide what to do
		return m, false, nil
	}

	if err != nil {
		return m, false, err
	}

	return m, true, nil
}

func (f finder) FindManyRaw(m interface{}, querySQL string, params ...interface{}) (interface{}, bool, error) {
	err := f.db.Raw(querySQL, params...).Find(m).Error

	if gorm.IsRecordNotFoundError(err) {
		// Let caller decide what to do
		return m, false, nil
	}

	if err != nil {
		return m, false, err
	}

	return m, true, nil
}
