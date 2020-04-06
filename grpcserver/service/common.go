package service

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	dbServiceName string = "grpc-db-service"
)

func InitDB(dbType string) {
	// Datadog
	sqltrace.Register(dbType,
		&mysql.MySQLDriver{},
		sqltrace.WithServiceName(dbServiceName),
	)
}

func OpenDB(deviceType, dataSource string) error {
	var err error
	db, err = sqltrace.Open(deviceType, dataSource)
	return err
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
