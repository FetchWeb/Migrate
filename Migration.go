package migrate

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Direction string

const (
	Up   Direction = "up"
	Down Direction = "down"
)

var DB *sql.DB

type Migration struct {
	Name        string
	Up          string
	Down        string
	IsInstalled bool
	InstalledAt time.Time
}

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

func (m *Migration) RunMigration(direction Direction) error {
	var migrationQueries []string

	migrationQueries = m.splitQueries(direction)

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

	if direction == Up {
		m.IsInstalled = true

		query := "INSERT INTO migration (name, up, down, is_installed) VALUES (?, ?, ?, ?);"
		_, err = transaction.Exec(query, m.Name, m.Up, m.Down, m.IsInstalled)
		if err != nil {
			if err := transaction.Rollback(); err != nil {
				return err
			}
			return err
		}
	} else {
		m.IsInstalled = false

		query := "DELETE FROM migration WHERE name = ?"

		_, err := transaction.Exec(query, m.Name)
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
