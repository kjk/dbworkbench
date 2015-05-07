package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sqlDb        *sql.DB
	errNYI             = errors.New("NYI")
	connectionID int32 = 1
	muCache      sync.Mutex
	dbUserCache  map[int]*User
)

// DbUser corresponds to users table
type DbUser struct {
	ID        int
	CreatedAt time.Time
	Email     string
	FullName  string
}

// User describes information about a user
type User struct {
	DbUser *DbUser
	// TODO: allow multiple connections
	ConnectionID int
	dbClient     *Client
}

func genNewConnectionID() int {
	id := atomic.AddInt32(&connectionID, 1)
	// TODO: wrap around after some reasonably large number
	return int(id)
}

func init() {
	dbUserCache = make(map[int]*User)
}

func isErrNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func dbGetUserByQuery(q string, args ...interface{}) (*DbUser, error) {
	var user DbUser
	db := getDbMust()
	err := db.QueryRow(q, args...).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.FullName)
	if isErrNoRows(err) {
		return nil, nil
	}
	if err != nil {
		LogInfof("db.QueryRow('%s') failed with %s\n", q, err)
		return nil, err
	}
	return &user, nil
}

func dbGetUserByID(id int) (*DbUser, error) {
	q := `SELECT id, created_at, email, full_name FROM users WHERE id=$1`
	return dbGetUserByQuery(q, id)
}

func dbGetUserByIDCached(id int) (*User, error) {
	LogInfof("id: %d\n", id)
	muCache.Lock()
	if dbUser, ok := dbUserCache[id]; ok {
		muCache.Unlock()
		return dbUser, nil
	}
	muCache.Unlock()

	dbUser, err := dbGetUserByID(id)
	if err != nil {
		return nil, err
	}
	user := &User{
		DbUser: dbUser,
	}
	muCache.Lock()
	dbUserCache[id] = user
	muCache.Unlock()
	return user, nil
}

func dbGetUserByEmail(email string) (*DbUser, error) {
	q := `SELECT id, created_at, email, full_name FROM users WHERE email=$1`
	return dbGetUserByQuery(q, email)
}

func dbGetOrCreateUser(email string, fullName string) (*DbUser, error) {
	LogInfof("email: %s, fullName: %s\n", email, fullName)
	dbUser, err := dbGetUserByEmail(email)
	if dbUser != nil {
		return dbUser, nil
	}
	if err != nil {
		// shouldn't happen, log and ignore
		LogErrorf("dbGetUserByEmail('%s') failed with '%s'\n", email, err)
	}

	db := getDbMust()
	q := `INSERT INTO users (email, fulL_name, created_at) VALUES ($1, $2, now())`
	_, err = db.Exec(q, email, fullName)
	if err != nil {
		return nil, err
	}
	return dbGetUserByEmail(email)
}

func getSqlConnectionRoot() string {
	if options.IsLocal {
		return "postgres://localhost/postgres?sslmode=disable"
	}
	// TODO: what is sslmode in prod?
	return "postgres://localhost/postgres"
}

func getSqlConnection() string {
	if options.IsLocal {
		return "postgres://localhost/dbworkbench?sslmode=disable"
	}
	// TODO: what is sslmode in prod?
	return "postgres://localhost/dbworkbench"
}

func execMust(db *sql.DB, q string, args ...interface{}) {
	LogVerbosef("db.Exec(): %s\n", q)
	_, err := db.Exec(q, args...)
	fatalIfErr(err, fmt.Sprintf("db.Exec('%s')", q))
}

func getCreateDbStatementsMust() []string {
	d, err := ioutil.ReadFile("createdb.sql")
	fatalIfErr(err, "getCreateDbStatementsMust")
	// can't execute multiple sql statements at once, so break the file
	// into separate statements
	return strings.Split(string(d), "\n\n")
}

func createDatabaseMust() *sql.DB {
	LogVerbosef("trying to create the database\n")
	db, err := sql.Open("postgres", getSqlConnectionRoot())
	fatalIfErr(err, "sql.Open()")
	LogVerbosef("got root connection\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	execMust(db, `CREATE DATABASE dbworkbench`)
	db.Close()

	db, err = sql.Open("postgres", getSqlConnection())
	fatalIfErr(err, "sql.Open()")
	// TODO: wrap in a transaction
	stmts := getCreateDbStatementsMust()
	for _, stm := range stmts {
		// skip empty lines
		stm = strings.TrimSpace(stm)
		if len(stm) > 0 {
			execMust(db, stm)
		}
	}

	LogVerbosef("created database\n")
	err = db.Ping()
	fatalIfErr(err, "db.Ping()")
	return db
}

func upgradeDbMust(db *sql.DB) {
	//q := `SELECT 1 FROM dbmigrations WHERE version = ?`
	//q := `INSERT INTO dbmigrations (version) VALUES (?)``
	// TODO: implement me
}

// note: no locking. the presumption is that this is called at startup and
// available throughout the lifetime of the program
func getDbMust() *sql.DB {
	if sqlDb != nil {
		return sqlDb
	}
	db, err := sql.Open("postgres", getSqlConnection())
	if err != nil {
		LogFatalf("sql.Open() failed with %s\n", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		if strings.Contains(err.Error(), `database "dbworkbench" does not exist`) {
			LogVerbosef("db.Ping() failed because no database exists\n")
			db = createDatabaseMust()
		} else {
			LogFatalf("db.Ping() failed with %s\n", err)
		}
	} else {
		upgradeDbMust(db)
	}
	sqlDb = db
	return sqlDb
}

func closeDb() {
	if sqlDb != nil {
		sqlDb.Close()
		sqlDb = nil
	}
}
