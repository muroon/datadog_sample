package service

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func OpenDB() error {
	var err error
	db, err = sql.Open("mysql", "root:@/grpc_datadog")
	return err
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

