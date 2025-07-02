package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type API struct {
	todos  TodoStore
	router *http.ServeMux
	addr   string
}

func (a *API) setup() {
	a.router.HandleFunc("GET /api/todos", a.listTodos)
	a.router.HandleFunc("GET /api/todos/{id}", a.getTodo)

	a.router.HandleFunc("POST /api/todos", a.createTodo)
	a.router.HandleFunc("DELETE /api/todos/{id}", a.deleteTodo)

	a.router.HandleFunc("PUT /api/todos/{id}", a.setCompleted)
}

func (a *API) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := a.todos.GetAll()
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

func (a *API) getTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := a.todos.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFound(w, r, err)
			return
		} else {
			a.internalServerError(w, r, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (a *API) createTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := readJson(w, r, &t); err != nil {
		a.badRequest(w, r, err)
		return
	}
	newID, err := a.todos.Create(t)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	t.ID = newID
	writeJSON(w, http.StatusOK, t)
}

func (a *API) setCompleted(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	todo, err := a.todos.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			a.notFound(w, r, err)
			return
		}
		a.internalServerError(w, r, err)
		return
	}

	if todo.Completed {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	todo.Completed = true
	if err := a.todos.Update(todo); err != nil {
		a.internalServerError(w, r, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.todos.Delete(id); err != nil {
		a.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *API) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s err: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "there was a problem")
}

func (app *API) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found: %s path: %s err: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}

func (app *API) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: %s path: %s err: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxByes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, map[string]string{"error": message})
}
