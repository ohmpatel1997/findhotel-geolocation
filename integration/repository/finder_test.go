package repository_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/stretchr/testify/assert"
)

type testFind struct {
	ID         int64
	SomeCol    string `gorm:"column:some_col"`
	AnotherCol string `gorm:"column:another_col"`
	CreatedAt  time.Time
}

type testFindUUID struct {
	ID         uuid.UUID
	SomeCol    string `gorm:"column:some_col"`
	AnotherCol string `gorm:"column:another_col"`
	CreatedAt  time.Time
}

func (tf *testFind) TableName() string {
	return "test_table_name"
}

func (tf *testFindUUID) TableName() string {
	return "test_table_name"
}

func TestFindByPKS(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("FindByPKS should find a row", func(t *testing.T) {
		tf := &testFind{}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(1, "foo", "bar")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."id" IN \(\$1\)\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs(1).
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTF, isFound, err := finder.FindByPKS(tf, 1)

		assert.Nil(err)

		assert.True(isFound)

		tf = foundTF.(*testFind)

		tf.CreatedAt = time.Time{}

		assert.Equal(&testFind{int64(1), "foo", "bar", time.Time{}}, tf)

	})

	t.Run("FindByPKS should find a uuid row", func(t *testing.T) {
		tf := &testFindUUID{}

		eID, err := uuid.NewUUID()

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow([]byte(eID.String()), "foo", "bar")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."id" IN \(\$1\)\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs(eID.String()).
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTF, isFound, err := finder.FindByPKS(tf, eID.String())

		assert.Nil(err)

		assert.True(isFound)

		tf = foundTF.(*testFindUUID)

		tf.CreatedAt = time.Time{}

		assert.Equal(&testFindUUID{eID, "foo", "bar", time.Time{}}, tf)

	})

	t.Run("FindBYPKS should not find a row", func(t *testing.T) {
		tf := &testFind{}

		// Return an empty row, found will be false
		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."id" IN \(\$1\)\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs(1).
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		_, isFound, err := finder.FindByPKS(tf, 1)

		assert.Nil(err)

		assert.False(isFound)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindBYPKS should not find a uuid row", func(t *testing.T) {
		tf := &testFindUUID{}

		uID := uuid.UUID{}

		// Return an empty row, found will be false
		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."id" IN \(\$1\)\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs(uID.String()).
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		_, isFound, err := finder.FindByPKS(tf, uID.String())

		assert.Nil(err)

		assert.False(isFound)

		if err := mock.ExpectationsWereMet(); err != nil {
			if !strings.Contains(err.Error(), "there is a remaining expectation which was not matched: ExpectedQuery => expecting Query, QueryContext or QueryRow which") {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}
	})

}

func TestFindFindManaged(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("FindManaged should find a row", func(t *testing.T) {
		tf := &testFind{
			SomeCol: "foo",
		}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(1, "foo", "bar")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."some_col" = \$1\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs("foo").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTF, found, err := finder.FindManaged(tf)

		tf = foundTF.(*testFind)

		assert.Nil(err)

		assert.True(found)

		tf.CreatedAt = time.Time{}

		assert.Equal(&testFind{int64(1), "foo", "bar", time.Time{}}, tf)
	})

	t.Run("FindManaged should not find a row", func(t *testing.T) {
		tf := &testFind{
			SomeCol: "foo",
		}

		// Return an empty row, found will be false
		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."some_col" = \$1\) ORDER BY "test_table_name"\."id" ASC LIMIT 1`).
			WithArgs("foo").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTF, found, err := finder.FindManaged(tf)

		tf = foundTF.(*testFind)

		assert.Nil(err)

		assert.False(found)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestFindManyManaged(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("FindManyManaged should find two rows", func(t *testing.T) {
		tfs := []testFind{}

		rows := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(1, "foo", "bar").
			AddRow(2, "foo", "bar2")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."some_col" = \$1\)`).
			WithArgs("foo").
			WillReturnRows(rows)

		finder := repository.NewFinder(gdb)

		foundTFS, err := finder.FindManyManaged(&tfs, repository.ColumnsAndValues{"some_col": "foo"})

		assert.Nil(err)

		tfs = *foundTFS.(*[]testFind)

		assert.Equal(2, len(tfs))

		tfs[0].CreatedAt = time.Time{}
		tfs[1].CreatedAt = time.Time{}

		assert.Equal(testFind{int64(1), "foo", "bar", time.Time{}}, tfs[0])
		assert.Equal(testFind{int64(2), "foo", "bar2", time.Time{}}, tfs[1])

	})

	t.Run("FindManyManaged should not return any rows", func(t *testing.T) {
		tfs := []testFind{}

		emptyRow := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE \("test_table_name"\."some_col" = \$1\)`).
			WithArgs("foo").
			WillReturnRows(emptyRow)

		finder := repository.NewFinder(gdb)

		foundTFS, err := finder.FindManyManaged(&tfs, repository.ColumnsAndValues{"some_col": "foo"})

		assert.Nil(err)

		tfs = *foundTFS.(*[]testFind)

		assert.Equal(0, len(tfs))
	})
}

func TestFindRaw(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("FindRaw should find a row", func(t *testing.T) {
		tf := &testFind{}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(1, "foo", "bar")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTF, isFound, err := finder.FindRaw(tf, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Nil(err)

		assert.True(isFound)

		tf = foundTF.(*testFind)

		tf.CreatedAt = time.Time{}

		assert.Equal(&testFind{int64(1), "foo", "bar", time.Time{}}, tf)
	})

	t.Run("FindRaw should not find a row", func(t *testing.T) {
		tf := &testFind{}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		_, isFound, err := finder.FindRaw(tf, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Nil(err)

		assert.False(isFound)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindRaw should error", func(t *testing.T) {
		tf := &testFind{}

		// Return an empty row, found will be false
		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(0, "foo", "bar").
			RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		_, _, err := finder.FindRaw(tf, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Equal("row error", err.Error())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestFindManyRaw(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("FindRaw should find a row", func(t *testing.T) {
		tfs := []testFind{}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(1, "foo", "bar").
			AddRow(2, "foo", "bar")

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTFS, isFound, err := finder.FindRaw(&tfs, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Nil(err)

		assert.True(isFound)

		tfs = *foundTFS.(*[]testFind)

		assert.Equal(2, len(tfs))

		assert.Equal(testFind{int64(1), "foo", "bar", time.Time{}}, tfs[0])
		assert.Equal(testFind{int64(2), "foo", "bar", time.Time{}}, tfs[1])
	})

	t.Run("FindRaw should not find a row", func(t *testing.T) {
		tfs := []testFind{}

		// Return an empty row, found will be false
		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"})

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTFS, _, err := finder.FindRaw(&tfs, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Nil(err)

		tfs = *foundTFS.(*[]testFind)

		assert.Equal(0, len(tfs))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindRaw should error", func(t *testing.T) {
		tfs := []testFind{}

		row := sqlmock.NewRows([]string{"id", "some_col", "another_col"}).
			AddRow(0, "foo", "bar").
			RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery(`SELECT \* FROM "test_table_name" WHERE some_col = \$1 AND another_col = \$2 ORDER BY some_col`).
			WithArgs("foo", "bar").
			WillReturnRows(row)

		finder := repository.NewFinder(gdb)

		foundTFS, _, err := finder.FindRaw(&tfs, `SELECT * FROM "test_table_name" WHERE some_col = ? AND another_col = ? ORDER BY some_col`, "foo", "bar")

		assert.Equal("row error", err.Error())

		tfs = *foundTFS.(*[]testFind)

		assert.Equal(0, len(tfs))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
