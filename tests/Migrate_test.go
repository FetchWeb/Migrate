package test

import (
	"testing"

	migrate "github.com/FetchWeb/Migrate"
)

func TestMigration(t *testing.T) {
	migrationOne := &migrate.Migration{}
	migrationOne.ParseSource("TestMigration_1.sql")

	migrationTwo := &migrate.Migration{}
	migrationTwo.ParseSource("TestMigration_2.sql")
}
