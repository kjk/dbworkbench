-- there should be at least one empty line between separate SQL statements
CREATE TABLE users (
  id                  INT NOT NULL AUTO_INCREMENT,
  created_at          TIMESTAMP NOT NULL,
  email               VARCHAR(255) NOT NULL,
  username            VARCHAR(255) NOT NULL,
  -- either password or google_oauth_json must be set
  password            VARCHAR(255),
  google_oauth_json   VARCHAR(2048),
  PRIMARY KEY (id),
  INDEX (email)
);

CREATE TABLE dbmigrations (
	version int NOT NULL
);
