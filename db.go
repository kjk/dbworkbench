package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

var (
	sqlDb  *sql.DB
	errNYI = errors.New("NYI")
)

var (
	// single muCache protects all things below
	muCache          sync.Mutex
	nextConnectionID = 1
	userCache        map[int]*User           // maps user id to User
	connections      map[int]*ConnectionInfo // maps connection id to ConnectionInfo
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
	ConnInfo *ConnectionInfo
}

// ConnectionInfo contains information about a database connection
type ConnectionInfo struct {
	ConnectionString string
	ConnectionID     int
	UserID           int
	CreatedAt        time.Time
	LastAccessAt     time.Time
	Client           *Client
}

func copyConnectionInfo(ci *ConnectionInfo) *ConnectionInfo {
	if ci == nil {
		return nil
	}
	return &ConnectionInfo{
		ConnectionString: ci.ConnectionString,
		ConnectionID:     ci.ConnectionID,
		UserID:           ci.UserID,
		CreatedAt:        ci.CreatedAt,
		LastAccessAt:     ci.LastAccessAt,
		Client:           ci.Client,
	}
}

func copyDbUser(u *DbUser) *DbUser {
	if u == nil {
		return nil
	}
	return &DbUser{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		Email:     u.Email,
		FullName:  u.FullName,
	}
}

// copyUser makes a deep copy of User so that it can be accessed without locking.
// Corollary: this is read-only, to make changes
func copyUser(u *User) *User {
	if u == nil {
		return nil
	}
	return &User{
		DbUser:   copyDbUser(u.DbUser),
		ConnInfo: copyConnectionInfo(u.ConnInfo),
	}
}

func getSqlConnectionRoot() string {
	if options.IsLocal {
		return "postgres://localhost/postgres?sslmode=disable"
	}
	return "postgres://dbworkbench_admin:f1r34nd1c3@127.0.0.1/postgres?sslmode=disable"
}

func getSqlConnection() string {
	if options.IsLocal {
		return "postgres://localhost/dbworkbench_admin?sslmode=disable"
	}
	return "postgres://dbworkbench_admin:f1r34nd1c3@127.0.0.1/dbworkbench_admin?sslmode=disable"
}

// creates new ConnectionInfo for a user and update user info. make sure
// to call removeCurrentUserConnectionInfo before calling this
func addConnectionInfo(connString string, userID int, client *Client) *ConnectionInfo {
	muCache.Lock()
	defer muCache.Unlock()
	nextConnectionID++
	conn := &ConnectionInfo{
		ConnectionString: connString,
		ConnectionID:     nextConnectionID,
		UserID:           userID,
		CreatedAt:        time.Now(),
		LastAccessAt:     time.Now(),
		Client:           client,
	}
	if user := userCache[userID]; user != nil {
		user.ConnInfo = conn
	} else {
		// TODO: shouldn't happen, log an error
	}
	return conn
}

func removeCurrentUserConnectionInfo(userID int) {
	muCache.Lock()
	defer muCache.Unlock()
	if user := userCache[userID]; user != nil {
		if user.ConnInfo != nil && user.ConnInfo.Client != nil {
			user.ConnInfo.Client.db.Close()
			user.ConnInfo = nil
		}
	} else {
		// TODO: shouldnt' happen, log an error
	}
}

func withUserLocked(userID int, f func(*User)) bool {
	foundUser := false
	muCache.Lock()
	if user := userCache[userID]; user != nil {
		f(user)
		foundUser = true
	}
	muCache.Unlock()
	return foundUser
}

func withConnectionLocked(connID int, f func(*ConnectionInfo)) bool {
	foundConn := false
	muCache.Lock()
	if conn := connections[connID]; conn != nil {
		f(conn)
		foundConn = true
	}
	muCache.Unlock()
	return foundConn
}

func updateConnectionLastAccess(connID int) {
	withConnectionLocked(connID, func(conn *ConnectionInfo) {
		conn.LastAccessAt = time.Now()
	})
}

func init() {
	userCache = make(map[int]*User)
	connections = make(map[int]*ConnectionInfo)
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

func dbGetUserCopyByIDCached(id int) (*User, error) {
	user, err := dbGetUserByIDCached(id)
	return copyUser(user), err
}

func dbGetUserByIDCached(id int) (*User, error) {
	//LogInfof("id: %d\n", id)
	muCache.Lock()
	if user, ok := userCache[id]; ok {
		muCache.Unlock()
		return user, nil
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
	userCache[id] = user
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
	q := `INSERT INTO users (email, full_name, created_at) VALUES ($1, $2, now())`
	_, err = db.Exec(q, email, fullName)
	if err != nil {
		return nil, err
	}
	return dbGetUserByEmail(email)
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
	execMust(db, `CREATE DATABASE dbworkbench_admin`)
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
		if strings.Contains(err.Error(), `database "dbworkbench_admin" does not exist`) {
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
