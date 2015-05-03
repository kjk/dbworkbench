-- there should be at least one empty line between separate SQL statements
CREATE EXTENSION hstore;

CREATE TABLE users (
  id                  SERIAL NOT NULL PRIMARY KEY,
  created_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  email               VARCHAR(255) NOT NULL,
  -- either password or google_oauth_json must be set
  password            VARCHAR(255),
  google_oauth_json   VARCHAR(2048)
);

CREATE INDEX idx_email ON users(email);

CREATE TABLE dbmigrations (
	version integer NOT NULL
);
