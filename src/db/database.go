package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Filter struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Record struct {
	Fields []Field `json:"fields"`
}

type Request struct {
	Table   string   `json:"table"`
	Fields  []Field  `json:"fields"`
	Filters []Filter `json:"filters"`
}

type Response struct {
	Records []Record `json:"records,omitempty"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

type Database struct {
	path string
	sql  *sql.DB
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("[Database] [Error] failed open database: %s", err)
	}

	db.Exec("DELETE FROM files")
	db.Exec("DELETE FROM filesystem")
	db.Exec("DELETE FROM tree_path")

	db.Exec("CREATE TABLE IF NOT EXISTS files (path TEXT PRIMARY KEY, data BLOB)")
	db.Exec("CREATE TABLE IF NOT EXISTS filesystem (id INTEGER PRIMARY KEY, path TEXT)")
	db.Exec("CREATE TABLE IF NOT EXISTS tree_path (ancestor INTEGER, descendant INTEGER, PRIMARY KEY (ancestor, descendant))")

	return &Database{
		path: path,
		sql:  db,
	}, nil
}

func (database *Database) Close() {
	database.sql.Close()
}
