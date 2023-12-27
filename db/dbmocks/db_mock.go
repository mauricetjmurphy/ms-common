package dbmocks

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSqlMock mocking db
func NewSqlMock() (*gorm.DB, sqlmock.Sqlmock) {
	var (
		sqlDb *sql.DB
		err   error
	)
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	if sqlDb == nil {
		panic("mock db is null")
	}
	conn, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic(err)
	}
	return conn, mock
}
