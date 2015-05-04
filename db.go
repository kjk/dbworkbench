package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	sqlDb   *sql.DB
	sqlDbMu sync.Mutex
	errNYI  = errors.New("NYI")
)

// DbUser corresponds to users table
type DbUser struct {
	ID    int
	Email string
}

func dbGetOrCreateUser(email string, fullName string) (*DbUser, error) {
	LogInfof("email: %s, fullName: %s\n", email, fullName)
	return nil, errNYI
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
	sqlDbMu.Lock()
	defer sqlDbMu.Unlock()
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
