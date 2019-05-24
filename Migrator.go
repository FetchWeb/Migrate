package migrate

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DB is the database where the migrations will be stored.
var DB *sql.DB

var MigrationDirectory string

func RunAllMigrations() error {
	migrationFiles := []os.FileInfo{}
	migrationEntries := []Migration{}

	files, err := ioutil.ReadDir(MigrationDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file)
		}
	}

	rows, err := DB.Query("SELECT * FROM migration WHERE deleted_at = NULL")
	if err != nil {
		return err
	}

	for rows.Next() {
		migrationEntry := Migration{}

		if err := rows.Scan(&migrationEntry); err != nil {
			return err
		}

		migrationEntries = append(migrationEntries, migrationEntry)
	}

	migrationsToRun := []Migration{}

	for _, migrationFile := range migrationFiles {
		migrationFound := false

		for _, migrationEntry := range migrationEntries {
			if strings.TrimSuffix(migrationFile.Name(), ".sql") == migrationEntry.Name {
				migrationFound = true
			}
		}

		if !migrationFound {
			newMigrationEntry := Migration{}
			newMigrationEntry.ParseSource(MigrationDirectory + migrationFile.Name())
			migrationsToRun = append(migrationsToRun, newMigrationEntry)
		}
	}

	for _, migrationToRun := range migrationsToRun {
		if err := migrationToRun.Run(Up); err != nil {
			return err
		}
	}

	return nil
}

func ListMigrations() ([]Migration, error) {
	migrationFiles := []os.FileInfo{}
	migrationEntries := []Migration{}

	files, err := ioutil.ReadDir(MigrationDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file)
		}
	}

	rows, err := DB.Query("SELECT * FROM migration WHERE deleted_at = NULL")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		migrationEntry := Migration{}

		if err := rows.Scan(&migrationEntry); err != nil {
			return nil, err
		}

		migrationEntries = append(migrationEntries, migrationEntry)
	}

	for _, migrationFile := range migrationFiles {
		migrationFound := false

		for _, migrationEntry := range migrationEntries {
			if strings.TrimSuffix(migrationFile.Name(), ".sql") == migrationEntry.Name {
				migrationFound = true
			}
		}

		if !migrationFound {
			newMigrationEntry := Migration{}
			newMigrationEntry.ParseSource(MigrationDirectory + migrationFile.Name())
			migrationEntries = append(migrationEntries, newMigrationEntry)
		}
	}

	return migrationEntries, nil
}
