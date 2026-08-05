package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {

	var err error

	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		dsn = "root:@tcp(localhost:3306)/absensi_mahasiswa?clientFoundRows=true"
	}

	DB, err = sql.Open(
		"mysql",
		dsn,
	)

	if err != nil {
		panic(err)
	}

	err = DB.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Database Connected")
}