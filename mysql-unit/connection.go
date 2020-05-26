package mysql_unit

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db = struct {
	host     string
	user     string
	driver   string
	password string
	database string
}{
	host:     "127.0.0.1:3305",
	user:     "root",
	driver:   "mysql",
	password: "root",
	database: "test_database",
}

func BasicConnection() *sql.DB {
	dsn := db.user + ":" + db.password + "@tcp(" + db.host + ")/" + db.database
	db, _ := sql.Open(db.driver, dsn)
	db.SetMaxIdleConns(0)

	return db
}
