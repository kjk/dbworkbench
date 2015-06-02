package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kjk/lzmadec"
	"github.com/kjk/u"
)

const (
	demoDBUser = "demodb"
)

var (
	dataDir string

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
  about_me            VARCHAR(32000),
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

CREATE TABLE badges (
	id 				SERIAL NOT NULL PRIMARY KEY,
	user_id 	INTEGER NOT NULL,
	name 			VARCHAR(256),
	date 			TIMESTAMP WITHOUT TIME ZONE
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

func init() {
	dataDir = u.ExpandTildeInPath("~/data/import_stack_overflow")
	os.MkdirAll(dataDir, 0755)
}

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
	name = fmt.Sprintf("%s.stackexchange.com.7z", name)
	return filepath.Join(dataDir, name)
}

// LogVerbosef logs in a verbose manner
func LogVerbosef(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	/*if pc, _, _, ok := runtime.Caller(1); ok {
		s = FunctionFromPc(pc) + ": " + s
	}*/
	fmt.Print(s)
}

// LogFatalf logs a message and quits
func LogFatalf(format string, arg ...interface{}) {
	LogVerbosef(format, arg)
	os.Exit(1)
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

	// grant demoDBUser read-only access to the database
	execMust(db, fmt.Sprintf(`GRANT USAGE ON SCHEMA public TO %s;`, demoDBUser))
	execMust(db, fmt.Sprintf(`GRANT SELECT ON users TO %s;`, demoDBUser))
	execMust(db, fmt.Sprintf(`GRANT SELECT ON posts TO %s;`, demoDBUser))
	execMust(db, fmt.Sprintf(`GRANT SELECT ON badges TO %s;`, demoDBUser))

	LogVerbosef("created database\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	return db
}

func httpDlAtomicCached(dstPath, uri string) error {
	// TODO: should at least check size of the file is correct
	if u.PathExists(dstPath) {
		LogVerbosef("'%s' already downloaded as '%s'\n", uri, dstPath)
		return nil
	}
	LogVerbosef("starting to download '%s'\n", uri)
	timeStart := time.Now()
	tmpPath := dstPath + ".tmp"
	fTmp, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer func() {
		if fTmp != nil {
			fTmp.Close()
			os.Remove(tmpPath)
		}
	}()

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	_, err = io.Copy(fTmp, resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	err = fTmp.Close()
	fTmp = nil
	if err != nil {
		return err
	}
	err = os.Rename(tmpPath, dstPath)
	if err != nil {
		return err
	}
	LogVerbosef("downloaded to '%s' in '%s'\n", dstPath, time.Since(timeStart))
	return nil
}

func getEntryForFile(archive *lzmadec.Archive, name string) *lzmadec.Entry {
	for _, e := range archive.Entries {
		if strings.EqualFold(name, e.Path) {
			return &e
		}
	}
	return nil
}

func toIntPtr(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

func toTimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func importSite(name string) error {
	db := createDatabaseMust(name)
	dstPath := getDataPath(name)
	uri := getDataURL(name)
	err := httpDlAtomicCached(dstPath, uri)
	if err != nil {
		return err
	}

	archive, err := lzmadec.NewArchive(dstPath)
	if err != nil {
		LogFatalf("lzmadec.NewArchive('%s') failed with '%s'\n", dstPath, err)
	}

	/*
		err = importPosts(archive, db)
		if err != nil {
			LogFatalf("importPosts() failed with %s\n", err)
		}

		err = importUsers(archive, db)
		if err != nil {
			LogFatalf("importUsers() failed with %s\n", err)
		}
	*/

	err = importBadges(archive, db)
	if err != nil {
		LogFatalf("importBadges() failed with %s\n", err)
	}
	return nil
}

func main() {
	err := importSite(sites[0])
	if err != nil {
		LogVerbosef("importSite() failed with %s\n", err)
	}
}
