package main

import (
	"database/sql"
	"strconv"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoStore struct {
	db *sql.DB
}

func (s *TodoStore) GetAll() ([]Todo, error) {
	rows, err := s.db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var (
			idInt     int64
			title     string
			completed bool
		)
		if err := rows.Scan(&idInt, &title, &completed); err != nil {
			return nil, err
		}
		todos = append(todos, Todo{
			ID:        strconv.FormatInt(idInt, 10),
			Title:     title,
			Completed: completed,
		})
	}
	return todos, nil
}

func (s *TodoStore) GetByID(id string) (Todo, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return Todo{}, err
	}
	var (
		title     string
		completed bool
	)
	err = s.db.QueryRow("SELECT title, completed FROM todos WHERE id = ?", idInt).Scan(&title, &completed)
	if err != nil {
		return Todo{}, err
	}
	return Todo{ID: id, Title: title, Completed: completed}, nil
}

func (s *TodoStore) Create(t Todo) (string, error) {
	res, err := s.db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", t.Title, t.Completed)
	if err != nil {
		return "", err
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(newID, 10), nil
}

func (s *TodoStore) Update(t Todo) error {
	idInt, err := strconv.ParseInt(t.ID, 10, 64)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?", t.Title, t.Completed, idInt)
	return err
}

func (s *TodoStore) Delete(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM todos WHERE id = ?", idInt)
	return err
}
