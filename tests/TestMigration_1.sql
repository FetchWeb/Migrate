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