package db

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

var fs embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(fs)
	goose.SetDialect("postgres")
	return goose.Up(db, ".")
}
