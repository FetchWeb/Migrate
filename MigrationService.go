package migrate

import (
	"database/sql"
)

type Direction string

const (
	Up   Direction = "up"
	Down Direction = "down"
)

var DB *sql.DB

func RunMigration(migration *Migration, direction Direction) error {
	var migrationQueries []string

	if direction == Up {
		migrationQueries = migration.Up
	} else {
		migrationQueries = migration.Down
	}

	transaction, err := DB.Begin()
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
