/* If users table exist, delete it */
DROP TABLE IF EXISTS users;

/* Create new table */
CREATE TABLE users (
	id VARCHAR(32) PRIMARY KEY,
	password VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
)