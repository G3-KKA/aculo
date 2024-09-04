package main

import (
	"database/sql"
	"embed"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	var db *sql.DB
	db = clickhouse.OpenDB(nil)
	// setup database

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect((string)(goose.DialectClickHouse)); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	// run app
}
