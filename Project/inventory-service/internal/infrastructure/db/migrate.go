package db

import (
	"database/sql"
	"embed"
	_ "embed"
	"github.com/pressly/goose/v3"
)

var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")
	return goose.Up(db, ".")
}
