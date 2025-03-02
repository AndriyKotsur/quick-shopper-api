package database

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migrate struct {
	DB *sql.DB
}

type Options func(opts *Migrate) error

func Migrator(db *sql.DB, opts ...Options) *Migrate {
	m := &Migrate{
		DB: db,
	}
	goose.SetBaseFS(embedMigrations)

	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return m
}

func (m *Migrate) Up() {
	if err := goose.Up(m.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (m *Migrate) Down() {
	if err := goose.Down(m.DB, "migrations"); err != nil {
		panic(err)
	}
}
