package models

import (
	"fmt"
	"log"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/jmoiron/sqlx"

	//
	// Allows sqlx to connect to Postgres.
	//
	_ "github.com/lib/pq"
)

var (
	Db *sqlx.DB
)

func init() {
	var err error
	Db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("dbname=%s sslmode=disable user=%s password=%s",
			os.Getenv("DATABASE"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
		))

	if err != nil {
		log.Fatal(err)
	}

	MigrateDB()
}

//
// MigrateDB loads all SQL files from migrations and executes them.
//
func MigrateDB() {
	migrationsBox, err := rice.FindBox("../migrations")

	if err != nil {
		log.Fatal(err)
	}

	migrationsBox.Walk("", func(path string, info os.FileInfo, err error) error {
		sql, _ := migrationsBox.String(path)
		Db.MustExec(sql)
		return nil
	})
}

//
// ResetDB resets the database schema, useful for testing.
//
func ResetDB() {
	Db.MustExec("drop schema public cascade")
	Db.MustExec("create schema public")
}
