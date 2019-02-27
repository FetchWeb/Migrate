package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

var db *sql.DB

func init() {
	dialect, isSet := os.LookupEnv("db_dialect")
	if !isSet {
		fmt.Printf("Failed to initialse Migrate package. Environment variable not set: db_dialect\n")
	}

	port, isSet := os.LookupEnv("db_port")
	if !isSet {
		fmt.Printf("Failed to initialse Migrate package. Environment variable not set: db_port\n")
	}

	username, isSet := os.LookupEnv("db_username")
	if !isSet {
		fmt.Printf("Failed to initialse Migrate package. Environment variable not set: db_username\n")
	}

	password, isSet := os.LookupEnv("db_password")
	if !isSet {
		fmt.Printf("Failed to initialse Migrate package. Environment variable not set: db_password\n")
	}

	name, isSet := os.LookupEnv("db_name")
	if !isSet {
		fmt.Printf("Failed to initialse Migrate package. Environment variable not set: db_name\n")
	}

	address := strings.Join([]string{
		username,
		":",
		password,
		"@tcp(127.0.0.1:",
		port,
		")/",
		name}, "")

	var err error
	db, err = sql.Open(dialect, address)
	if err != nil {
		fmt.Printf("Failed to open Migrate database connection: %v", err)
	}
}

func RunMigration(migrationQueries []string) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	for _, query := range migrationQueries {
		_, err := transaction.Exec(query)

		if err != nil {
			if err := transaction.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}
