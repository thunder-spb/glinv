package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"github.com/vharitonsky/iniflags"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"glinv/pkg/models/pgsql"
)

type contextKey string

var contextKeyUser = contextKey("user")

// Define an application struct to hold the application-wide
// dependencies for the web application.
type application struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	session          *sessions.Session
	users            *pgsql.UserModel
	inventoryHost    *pgsql.InventoryHostModel
	inventoryHVar    *pgsql.InventoryHVarModel
	inventoryHTag    *pgsql.InventoryHTagModel
	inventoryGroup   *pgsql.InventoryGroupModel
	inventoryGVar    *pgsql.InventoryGVarModel
	inventoryService *pgsql.InventoryServiceModel
	serverAgent      *pgsql.ServerAgentModel
	baseTemplate     *pgsql.BaseTemplateModel
	history          *pgsql.HistoryModel
	templateCache    map[string]*template.Template
}

func main() {
	//conf := LoadConfiguration("vault/secrets/pgsql-glinv-creds")
	migrationDir := flag.String("migrationDir", "./db/migrations", "Directory where the migration files are located ?")
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:password@192.168.10.229:5432/glinv?sslmode=disable", "PostgreSQL data source name")
	//dsn := flag.String("dsn", "postgres://postgres:password@localhost:5432/glinv?sslmode=disable", "PostgreSQL data source name")

	// For encrypt and authenticate session cookies, it should be 32 bytes long.
	secret := flag.String("secret", "u46ImCV9y5Vlur8YvODJEhgOY8p9JVE4", "Secret key")
	// Sessions always expires after 12 hours.
	sessionTime := flag.Int("sessiontime", 12, "Session lifetime in hours")
	// session.SameSite = http.SameSiteStrictMode // TODO

	//flag.Parse()
	iniflags.Parse() // use instead of flag.Parse()c

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager, passing in the secret key as the parameter.
	session := sessions.New([]byte(*secret))
	session.Lifetime = time.Duration(*sessionTime) * time.Hour

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Run migrations
	migration(db, migrationDir)

	go sJobs(db)
	go mJobs(db)
	go hJobs(db)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		session:          session,
		users:            &pgsql.UserModel{DB: db},
		inventoryHost:    &pgsql.InventoryHostModel{DB: db},
		inventoryHVar:    &pgsql.InventoryHVarModel{DB: db},
		inventoryHTag:    &pgsql.InventoryHTagModel{DB: db},
		inventoryGroup:   &pgsql.InventoryGroupModel{DB: db},
		inventoryGVar:    &pgsql.InventoryGVarModel{DB: db},
		inventoryService: &pgsql.InventoryServiceModel{DB: db},
		serverAgent:      &pgsql.ServerAgentModel{DB: db},
		baseTemplate:     &pgsql.BaseTemplateModel{DB: db},
		history:          &pgsql.HistoryModel{DB: db},
		templateCache:    templateCache,
	}

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and
// returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Database migrations
func migration(db *sql.DB, migrationDir *string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir), // file://path/to/directory
		"postgres", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")
	//os.Exit(0)
}
