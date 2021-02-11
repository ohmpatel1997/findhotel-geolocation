package repository_test

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/stretchr/testify/assert"
)

type testModel struct {
	ID        int64
	SomeCol   string `gorm:"column:some_col"`
	CreatedAt time.Time
}

func (tm *testModel) TableName() string {
	return "test_table_name"
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		SomeCol: "foo",
	}

	row := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	// GORM + postgres uses query instead of Exec to INSERT.. because of reasons?
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT  INTO "test_table_name" \("some_col","created_at"\) VALUES \(\$1\,\$2\) RETURNING "test_table_name"."id"`).
		WithArgs("foo", AnyTime{}).
		WillReturnRows(row)
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Insert(tm)

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertWithTransaction(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		SomeCol: "foo",
	}

	row := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	// GORM + postgres uses query instead of Exec to INSERT.. because of reasons?
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT  INTO "test_table_name" \("some_col","created_at"\) VALUES \(\$1\,\$2\) RETURNING "test_table_name"."id"`).
		WithArgs("foo", AnyTime{}).
		WillReturnRows(row)
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Transact(func(tx *gorm.DB) error {
		return cuder.InsertWithTX(tx, tm)
	})

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		ID:        1,
		SomeCol:   "foo",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "test_table_name" SET "some_col" = \$1  WHERE "test_table_name"\."id" = \$2`).
		WithArgs("new_val", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Update(tm, repository.ColumnsAndValues{"some_col": "new_val"})

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateWithTransaction(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		ID:        1,
		SomeCol:   "foo",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "test_table_name" SET "some_col" = \$1  WHERE "test_table_name"\."id" = \$2`).
		WithArgs("new_val", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Transact(func(tx *gorm.DB) error {
		return cuder.UpdateWithTX(tx, tm, repository.ColumnsAndValues{"some_col": "new_val"})
	})

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		ID:        1,
		SomeCol:   "foo",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "test_table_name" WHERE "test_table_name"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Delete(tm)

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteWithTransaction(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tm := &testModel{
		ID:        1,
		SomeCol:   "foo",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "test_table_name" WHERE "test_table_name"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	cuder := repository.NewCuder(gdb)

	err = cuder.Transact(func(tx *gorm.DB) error {
		return cuder.DeleteWithTX(tx, tm)
	})

	assert.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
