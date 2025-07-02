package main

import (
	"database/sql"
)

type Todo struct {
	id        string
	title     string
	completed bool
}

type TodoStore struct {
	db *sql.DB
}
