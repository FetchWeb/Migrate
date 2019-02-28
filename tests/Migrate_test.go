package test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	migrate "github.com/FetchWeb/Migrate"
)

func TestMigration(t *testing.T) {
	config := &struct {
		DBDialect  string `json:"db_dialect"`
		DBPort     string `json:"db_port"`
		DBUsername string `json:"db_username"`
		DBPassword string `json:"db_password"`
		DBName     string `json:"db_name"`
	}{}

	configFile, err := ioutil.ReadFile("TestConfig.json")
	if err != nil {
		t.Fatalf("Failed to read TestConfig file: %v", err)
	}

	err = json.Unmarshal(configFile, config)
	if err != nil {
		t.Fatalf("Failed to unmarshal TestConfig file: %v", err)
	}

	address := strings.Join([]string{
		config.DBUsername,
		":",
		config.DBPassword,
		"@tcp(127.0.0.1:",
		config.DBPort,
		")/"}, "")

	migrate.DB, err = sql.Open(config.DBDialect, address)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	migrate.DB.Exec(strings.Join([]string{"DROP DATABASE ", config.DBName, ";"}, ""))

	_, err = migrate.DB.Exec(strings.Join([]string{"CREATE DATABASE ", config.DBName, ";"}, ""))
	if err != nil {
		t.Fatalf("Failed to create %v database: %v", config.DBName, err)
	}

	_, err = migrate.DB.Exec(strings.Join([]string{"USE ", config.DBName, ";"}, ""))
	if err != nil {
		t.Fatalf("Failed to use %v database: %v", config.DBName, err)
	}

	migrationOne := &migrate.Migration{}
	if err := migrationOne.ParseSource("TestMigration_1.sql"); err != nil {
		t.Fatalf("Failed to parse migrationOne source: %v", err)
	}

	if err := migrationOne.RunMigration(migrate.Up); err != nil {
		t.Fatalf("Failed to migrate up on migrationOne: %v", err)
	}

	if err := migrationOne.RunMigration(migrate.Down); err != nil {
		t.Fatalf("Failed to migrate down on migrationOne: %v", err)
	}

	migrationTwo := &migrate.Migration{}
	if err := migrationTwo.ParseSource("TestMigration_2.sql"); err != nil {
		t.Fatalf("Failed to parse migrationTwo source: %v", err)
	}

	if err := migrationTwo.RunMigration(migrate.Up); err != nil {
		t.Fatalf("Failed to migrate up on migrationTwo: %v", err)
	}

	_, err = migrate.DB.Exec(strings.Join([]string{"DROP DATABASE ", config.DBName, ";"}, ""))
	if err != nil {
		t.Fatalf("Failed to remove %v database: %v", config.DBName, err)
	}
}
