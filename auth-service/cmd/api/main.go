package main

import (
	"auth-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Setup Config
	app := &Config{
		Client: &http.Client{},
	}
	app.setupRepository(conn)
	log.Printf("Starting Auth Service on PORT %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts += 1
		} else {
			log.Println("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepository(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
