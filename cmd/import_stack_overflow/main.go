package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kjk/lzmadec"
	"github.com/kjk/stackoverflow"
	"github.com/kjk/u"
	"github.com/lib/pq"
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

	LogVerbosef("created database\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	return db
}

func httpDlAtomicCached(dstPath, uri string) error {
	// TODO: should at least check size of the file is correct
	if u.PathExists(dstPath) {
		LogVerbosef("'%s' already download as '%s'\n", uri, dstPath)
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

func importUsersIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			LogVerbosef("calling txn.Rollback()\n")
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("users",
		"id",
		"reputation",
		"creation_date",
		"display_name",
		"last_access_date",
		"website_url",
		"location",
		"about_me",
		"views",
		"up_votes",
		"down_votes",
		"age",
		"account_id",
		"profile_image_url"))
	if err != nil {
		LogVerbosef("txn.Prepare() failed with %s\n", err)
		return 0, fmt.Errorf("txt.Prepare() failed with %s", err)
	}
	n := 0
	for r.Next() {
		u := &r.User
		_, err = stmt.Exec(
			u.ID,
			u.Reputation,
			u.CreationDate,
			toStringPtr(u.DisplayName),
			toTimePtr(u.LastAccessDate),
			toStringPtr(u.WebsiteURL),
			toStringPtr(u.Location),
			toStringPtr(u.AboutMe),
			u.Views,
			u.UpVotes,
			u.DownVotes,
			u.Age,
			u.AccountID,
			toStringPtr(u.ProfileImageURL),
		)
		if err != nil {
			LogVerbosef("stmt.Exec() failed with %s\n", err)
			fmt.Printf("len(u.DisplayName): %d\n", len(u.DisplayName))
			fmt.Printf("len(u.WebsiteURL): %d\n", len(u.WebsiteURL))
			fmt.Printf("len(u.Location): %d\n", len(u.Location))
			fmt.Printf("len(u.AboutMe): %d\n", len(u.AboutMe))
			fmt.Printf("u.AboutMe: '%s'\n", u.AboutMe)
			return 0, fmt.Errorf("stmt.Exec() failed with %s", err)
		}
		n++
	}
	if err = r.Err(); err != nil {
		LogVerbosef("r.Err() failed with %s\n", err)
		return 0, err
	}
	_, err = stmt.Exec()
	if err != nil {
		LogVerbosef("stmt.Exec() 2 failed with %s\n", err)
		return 0, fmt.Errorf("stmt.Exec() failed with %s", err)
	}
	err = stmt.Close()
	if err != nil {
		LogVerbosef("stmt.Close() failed with %s\n", err)
		return 0, fmt.Errorf("stmt.Close() failed with %s", err)
	}
	err = txn.Commit()
	txn = nil
	if err != nil {
		LogVerbosef("txn.Commit() failed with %s\n", err)
		return 0, fmt.Errorf("txn.Commit() failed with %s", err)
	}
	return n, nil
}

func importUsers(archive *lzmadec.Archive, db *sql.DB) error {
	usersRecord := getEntryForFile(archive, "Users.xml")
	if usersRecord == nil {
		return errors.New("genEntryForFile('Users.xml') returned nil")
	}

	usersReader, err := archive.ExtractReader(usersRecord.Path)
	if err != nil {
		return fmt.Errorf("ExtractReader('%s') failed with %s", usersRecord.Path, err)
	}
	defer usersReader.Close()
	ur, err := stackoverflow.NewUsersReader(usersReader)
	if err != nil {
		return fmt.Errorf("stackoverflow.NewUsersReader() failed with %s", err)
	}
	n, err := importUsersIntoDB(ur, db)
	if err != nil {
		return fmt.Errorf("importUsersIntoDB() failed with %s", err)
	}
	LogVerbosef("processed %d user records\n", n)
	return nil
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

	err = importUsers(archive, db)
	if err != nil {
		LogFatalf("importUsers() failed with %s\n", err)
	}
	return nil
}

func main() {
	importSite(sites[0])
}
