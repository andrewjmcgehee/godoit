package main

import (
	"database/sql"
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/andrewjmcgehee/godo/internal/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Database struct {
	db      *sql.DB
	Queries *orm.Queries
}

func NewDatabase() (*Database, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	localShareGodo := filepath.Join(home, ".local", "share", "godo")
	err = os.MkdirAll(localShareGodo, 0755)
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(localShareGodo, "todos.db")
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// disable goose logs
	logWriter := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(logWriter)

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("failed to migrate todos.db")
	}
	return &Database{
		db:      sqlDB,
		Queries: orm.New(sqlDB),
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

type Priority string

const (
	P0 Priority = "P0"
	P1 Priority = "P1"
	P2 Priority = "P2"
)

type Todo struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Priority  Priority  `json:"priority"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
