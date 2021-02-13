package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	SSLModeRequire = "require"

	SSLModeDisable = "disable"

	pgDSNTemplate = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"

	defaultSetMaxOpenConns = 1
	defaultSetMaxIdleConns = 1
)

type PGConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string

	//Optional
	SetMaxOpenConns    *int
	SetMaxIdleConns    *int
	SetConnMaxLifetime *time.Duration
}

func NewPGConnection(pgc *PGConfig, connStr *string) (*gorm.DB, error) {
	var err error = nil
	dsn := ""

	if connStr != nil {
		dsn = *connStr
	} else if pgc != nil {
		dsn, err = pgc.String()
	} else {
		return nil, fmt.Errorf("please specify DB configuration")
	}

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if pgc != nil {
		if pgc.SetMaxOpenConns != nil {
			db.DB().SetMaxOpenConns(*pgc.SetMaxOpenConns)
		} else {
			db.DB().SetMaxOpenConns(defaultSetMaxOpenConns)
		}

		if pgc.SetMaxIdleConns != nil {
			db.DB().SetMaxIdleConns(*pgc.SetMaxIdleConns)
		} else {
			db.DB().SetMaxIdleConns(defaultSetMaxIdleConns)
		}

		//Default for this is a lifetime conn
		if pgc.SetConnMaxLifetime != nil {
			db.DB().SetConnMaxLifetime(*pgc.SetConnMaxLifetime)
		}
	}
	err = db.DB().Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (pgc PGConfig) String() (string, error) {
	var s string

	if pgc.Host == "" {
		return s, errors.New("Host must be set")
	}

	if pgc.Port == 0 {
		return s, errors.New("Port must be set")
	}

	if pgc.User == "" {
		return s, errors.New("User must be set")
	}

	if pgc.Password == "" {
		return s, errors.New("Password must be set")
	}

	if pgc.DBName == "" {
		return s, errors.New("DBName must be set")
	}

	if pgc.SSLMode == "" {
		return s, errors.New("SSLMode must be set")
	}

	s = fmt.Sprintf(pgDSNTemplate, pgc.Host, pgc.Port, pgc.User, pgc.Password, pgc.DBName, pgc.SSLMode)

	return s, nil
}
