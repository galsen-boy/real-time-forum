package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//Opens the database for sql execution and querying
func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//Opens the database and executes the initialisation statement, creating the database if it does not exist
func InitDB(path string) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(CreateTables)
	if err != nil {
		return err
	}

	return nil
}