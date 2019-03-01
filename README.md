# FetchWeb Migrate

## Introduction
FetchWeb Migrate is a simple database migration API witten in Go with no dependencies outside of the standard library. Migrations offer a way of reliably rolling out and rolling back database changes.

## Setup Example

### 1. Create Migration source file

Migration source files are written in SQL and must up the "-- UP" and "-- DOWN" comments in them to that the Migration can split the source file when it is parsed.

```sql
-- UP
CREATE TABLE migration (
	id BIGINT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	up TEXT NOT NULL,
	down TEXT,
	is_installed BIT DEFAULT 0,
	installed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	deleted_at DATETIME,
	PRIMARY KEY (id)
);

CREATE TABLE test_table (
	id BIGINT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	deleted_at DATETIME,
	PRIMARY KEY (id)
);

INSERT INTO test_table (name) VALUES ('test_name');

-- DOWN
TRUNCATE TABLE test_table;

DROP TABLE test_table;
```

### 2. Implement Migrate API

Create a Migration then parse the SQL source file to it.

From there a Migration can be run in two directions, Up or Down. Running a Migration Up will rollout database changes and add the Migration info to the database, so a Migration table is required in the database, see the test Migration sources for the insert query. Running a Migration Down will rollback database changes and set the IsInstalled column to false for that Migration in the database.

Migration names should be unque to avoid confusion. Prepending a UNIX timestamp or the version that this Migration corresponds with can help ensure they are. For example.

```
1551438141_NewMigration.sql
150_NewMinorVersion.sql
```

```go
package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	migrate "github.com/FetchWeb/Migrate"
)

func main() {
	// Open connection to the database.
	migrate.DB, err = sql.Open("mysql", "<username>:<password>@tcp(127.0.0.1:<port>)/<name>")
	if err != nil {
		panic(err)
	}
	defer migrate.DB.close()

	// Create Migration and parse the source file.
	migrationOne := &migrate.Migration{}
	if err := migrationOne.ParseSource("<file name>.sql"); err != nil {
		panic(err)
	}

	// Run the Migration Up.
	if err := migrationOne.RunMigration(migrate.Up); err != nil {
		panic(err)
	}

	// Run the Migration Down.
	if err := migrationOne.RunMigration(migrate.Down); err != nil {
		panic(err)
	}
}
```