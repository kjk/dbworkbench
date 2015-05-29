package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kjk/u"
	_ "github.com/lib/pq"
)

const (
	demoDBUser = "demodb"
)

var (
	sites = []string{
		"academia",
		"android",
		"anime",
	}

	initialSchema = `
CREATE TABLE users (
  id                  SERIAL NOT NULL PRIMARY KEY,
  reputation          INTEGER,
  creation_date       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  display_name        VARCHAR(255),
  last_access_date    TIMESTAMP WITHOUT TIME ZONE,
  website_url         VARCHAR(512),
  location            VARCHAR(1024),
  about_me            VARCHAR(4096),
  views               INTEGER NOT NULL,
  up_votes            INTEGER NOT NULL,
  down_votes          INTEGER NOT NULL,
  age                 INTEGER NOT NULL,
  account_id          INTEGER NOT NULL,
  profile_image_url   VARCHAR(512)
);

CREATE TABLE posts (
	id                       SERIAL NOT NULL PRIMARY KEY,
  post_type_id             INTEGER NOT NULL,
  parent_id                INTEGER,
  accepted_answer_id       INTEGER,
  creation_date            TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  score                    INTEGER NOT NULL,
  view_count               INTEGER NOT NULL,
  body                     TEXT NOT NULL,
  owner_user_id            INTEGER,
  owner_display_name       VARCHAR(512),
  last_editor_user_id      INTEGER,
  last_editor_display_name VARCHAR(512),
  last_edit_date           TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  last_activity_date       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  title                    VARCHAR(512),
  tags                     VARCHAR(2048),
  answer_count             INTEGER NOT NULL,
  comment_count            INTEGER NOT NULL,
  favorite_count           INTEGER NOT NULL,
  community_owned_date     TIMESTAMP WITHOUT TIME ZONE,
  closed_date              TIMESTAMP WITHOUT TIME ZONE
);
`

	// http://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
	createDemodbRoleIfNotExistsStmt = `
do
$body$
declare
  num_users integer;
begin
   SELECT count(*)
     into num_users
   FROM pg_user
   WHERE usename = 'demodb';

   IF num_users = 0 THEN
      CREATE ROLE demodb LOGIN PASSWORD 'demodb';
   END IF;
end
$body$
;`
)

func getSqlConnectionRoot() string {
	return "postgres://localhost/postgres?sslmode=disable"
}

func getSqlConnectionForDB(name string) string {
	return fmt.Sprintf("postgres://localhost/%s?sslmode=disable", name)
}

func getDataURL(name string) string {
	return fmt.Sprintf("https://archive.org/download/stackexchange/%s.stackexchange.com.7z", name)
}

func getDataPath(name string) string {
	return fmt.Sprintf("%s.stackexchange.com.7z", name)
}

func LogVerbosef(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	/*if pc, _, _, ok := runtime.Caller(1); ok {
		s = FunctionFromPc(pc) + ": " + s
	}*/
	fmt.Print(s)
}

func fatalIfErr(err error, what string) {
	if err != nil {
		log.Fatalf("%s failed with %s\n", what, err)
	}
}

func execMust(db *sql.DB, q string, args ...interface{}) {
	LogVerbosef("db.Exec(): %s\n", q)
	_, err := db.Exec(q, args...)
	fatalIfErr(err, fmt.Sprintf("db.Exec('%s')", q))
}

func getCreateDbStatementsMust() []string {
	// can't execute multiple sql statements at once, so break the file
	// into separate statements
	return strings.Split(initialSchema, "\n\n")
}

func createDatabaseMust(dbName string) *sql.DB {
	LogVerbosef("trying to create the database\n")
	db, err := sql.Open("postgres", getSqlConnectionRoot())
	fatalIfErr(err, "sql.Open()")
	LogVerbosef("got root connection\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	execMust(db, createDemodbRoleIfNotExistsStmt)
	execMust(db, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	execMust(db, fmt.Sprintf("CREATE DATABASE %s", dbName))
	execMust(db, fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", dbName, demoDBUser))
	db.Close()

	db, err = sql.Open("postgres", getSqlConnectionForDB(dbName))
	fatalIfErr(err, "sql.Open()")
	stmts := getCreateDbStatementsMust()
	for _, stm := range stmts {
		// skip empty lines
		stm = strings.TrimSpace(stm)
		if len(stm) > 0 {
			execMust(db, stm)
		}
	}

	// grant dmoDBUser read-only access to the database
	execMust(db, fmt.Sprintf(`GRANT USAGE ON SCHEMA public TO %s;`, demoDBUser))
	execMust(db, fmt.Sprintf(`GRANT SELECT ON %s TO %s;`, dbName, demoDBUser))

	LogVerbosef("created database\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	return db
}

func httpDlAtomicCached(dstPath, uri string) error {
	// TODO: should at least check size of the file is correct
	if u.PathExists(dstPath) {
		return nil
	}
	tmpPath := dstPath + ".tmp"
	fTmp, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer func() {
		if fTmp != nil {
			fTmp.Close()
		}
	}()

	return nil
}

func importSite(name string) error {
	//db := createDatabaseMust(name)
	dstPath := getDataPath(name)
	uri := getDataURL(name)
	err := httpDlAtomicCached(dstPath, uri)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	importSite(sites[0])
}
