package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbAddr := envGetString("DB_ADDR", ":memory:")
	db, err := newDB(dbAddr, 10, 5, "30m")
	if err != nil {
		log.Fatal(err)
	}
	app := API{
		todos:  TodoStore{db},
		router: http.NewServeMux(),
		addr:   envGetString("ADDR", ":8080"),
	}

	app.setup()

	log.Printf("starting server at %s", app.addr)
	if err := http.ListenAndServe(app.addr, app.router); err != nil {
		log.Fatal(err)
	}

}

func newDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("mysql", addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func envGetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
