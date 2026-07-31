package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {

	var err error

	DB, err = sql.Open(
		"mysql",
		"root:@tcp(localhost:3306)/absensi_mahasiswa?clientFoundRows=true",
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