package test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	migrate "github.com/FetchWeb/Migrate"
)

func TestListMigrations(t *testing.T) {
	config := &struct {
		DBDialect  string `json:"db_dialect"`
		DBPort     string `json:"db_port"`
		DBUsername string `json:"db_username"`
		DBPassword string `json:"db_password"`
		DBName     string `json:"db_name"`
	}{}

	// Read the test config and unmarshal.
	configFile, err := ioutil.ReadFile("TestConfig.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if err := json.Unmarshal(configFile, config); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Open connection to the database.
	migrate.DB, err = sql.Open(config.DBDialect, config.DBUsername+":"+config.DBPassword+"@tcp(127.0.0.1:"+config.DBPort+")/")
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}

	// Drop any remaining test databases that might have been left from previous tests.
	migrate.DB.Exec(strings.Join([]string{"DROP DATABASE ", config.DBName, ";"}, ""))

	// Create and select the database for running the test.
	_, err = migrate.DB.Exec(strings.Join([]string{"CREATE DATABASE ", config.DBName, ";"}, ""))
	if err != nil {
		t.Fatalf("Failed to create %v database: %v", config.DBName, err)
	}
	_, err = migrate.DB.Exec(strings.Join([]string{"USE ", config.DBName, ";"}, ""))
	if err != nil {
		t.Fatalf("Failed to use %v database: %v", config.DBName, err)
	}
	_, err = migrate.DB.Exec("CREATE TABLE migration (id BIGINT NOT NULL AUTO_INCREMENT, name VARCHAR(255) NOT NULL, up TEXT NOT NULL, down TEXT, is_installed BIT DEFAULT 0, installed_at DATETIME DEFAULT CURRENT_TIMESTAMP, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, deleted_at DATETIME, PRIMARY KEY (id))")
	if err != nil {
		t.Fatalf("Failed to use %v database: %v", config.DBName, err)
	}

	migrate.MigrationDirectories = []string{"/home/strongishllama/go/src/github.com/FetchWeb/Migrate/tests/migrations"}
	migrations, err := migrate.ListMigrations()
	if err != nil {
		t.Fatalf("Failed to list migrations: %v", err)
	}

	t.Log(migrations)
}
