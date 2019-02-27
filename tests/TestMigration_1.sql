-- UP
CREATE DATABASE migration_test;

CREATE TABLE test_table (
	id BIGINT NOT NULL AUTO_INCREMENT,
	title VARCHAR(255) NOT NULL,
	created_at DATETIME SET DEFAULT GETDATE(),
	deleted_at DATETIME
	PRIMARY KEY (id)
);

-- DOWN
DROP TABLE test_table;

DROP DATABASE migration_test;