-- UP
CREATE TABLE test_table (
	id BIGINT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	deleted_at DATETIME,
	PRIMARY KEY (id)
);

INSERT INTO test_table (name) VALUES ('test_name');