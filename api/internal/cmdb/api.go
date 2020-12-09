package cmdb

import (
	"database/sql"
	"log"
	"net/http"
	"restapi/internal/store/pgsql"
	"time"

	_ "github.com/lib/pq" // ...
)

// Start API
func Start(config *Config) error {
	db, err := openDB(config.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	store := pgsql.New(db)
	srv := newServer(store)

	log.Println("Starting CMDB REST API on", config.Addr)
	return http.ListenAndServe(config.Addr, srv)
}

// Open DB
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
