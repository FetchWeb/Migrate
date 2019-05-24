package migrate

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// Direction type alias.
type Direction string

// The possible directions a Migration can be run.
const (
	Up   Direction = "up"
	Down Direction = "down"
)

// Migration gives the ability to reliably rollout and rollback database changes.
type Migration struct {
	ID          int64
	Name        string
	Up          string
	Down        string
	IsInstalled bool
	InstalledAt time.Time
}

// ParseSource parses the SQL source file into the Migration.
func (m *Migration) ParseSource(fileDirectory string) error {
	m.Name = strings.TrimSuffix(filepath.Base(fileDirectory), ".sql")

	file, err := ioutil.ReadFile(fileDirectory)
	if err != nil {
		return err
	}

	queries := strings.Split(string(file), "-- DOWN")

	m.Up = strings.TrimLeft(queries[0], "-- UP")

	if len(queries) > 1 {
		m.Down = strings.TrimSpace(queries[1])
	}

	return nil
}

// Run takes a direction as an argument and runs the queries for that direction.
func (m *Migration) Run(direction Direction) error {
	if direction != Up && direction != Down {
		return errors.New("Invalid Migration Direction")
	}

	// Split up the queries.
	var migrationQueries []string
	migrationQueries = m.splitQueries(direction)

	// Start a transaction so a rollback can happen if a query fails.
	transaction, err := DB.Begin()
	if err != nil {
		return err
	}

	// Execute each query.
	for _, query := range migrationQueries {
		_, err := transaction.Exec(query)

		if err != nil {
			if err := transaction.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	// If the Direction was up, insert the migration into the database.
	// If the Direction was down, set the is_installed column to false.
	if direction == Up {
		m.IsInstalled = true

		query := "INSERT INTO migration (name, up, down, is_installed) VALUES (?, ?, ?, ?);"
		res, err := transaction.Exec(query, m.Name, m.Up, m.Down, m.IsInstalled)
		if err != nil {
			if err := transaction.Rollback(); err != nil {
				return err
			}
			return err
		}

		m.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		m.IsInstalled = false

		query := "UPDATE migration SET is_installed = 0 WHERE id = ?"

		_, err := transaction.Exec(query, m.ID)
		if err != nil {
			if err := transaction.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	// Commit the transaction.
	return transaction.Commit()
}

// splitQueries splits a Migration's Up or Down into separate queries.
func (m *Migration) splitQueries(direction Direction) []string {
	var queries []string

	if direction == Up {
		queries = strings.SplitAfter(m.Up, ";")
		queries = queries[:len(queries)-1]

		for i := 0; i < len(queries); i++ {
			queries[i] = strings.TrimSpace(queries[i])
		}
	} else {
		if m.Down != "" {
			queries = strings.SplitAfter(m.Down, ";")
			queries = queries[:len(queries)-1]

			for i := 0; i < len(queries); i++ {
				queries[i] = strings.TrimSpace(queries[i])
			}
		}
	}

	return queries
}
