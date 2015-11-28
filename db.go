package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	errNYI = errors.New("NYI")
)

var (
	// single muCache protects all things below
	muCache          sync.Mutex
	nextConnectionID = 1
	connections      map[int]*ConnectionInfo // maps connection id to ConnectionInfo
)

// ConnectionInfo contains information about a database connection
type ConnectionInfo struct {
	ConnectionString string
	ConnectionID     int
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
		CreatedAt:        ci.CreatedAt,
		LastAccessAt:     ci.LastAccessAt,
		Client:           ci.Client,
	}
}

// creates new ConnectionInfo for a user and update user info. make sure
// to call removeCurrentUserConnectionInfo before calling this
func addConnectionInfo(connString string, client *Client) *ConnectionInfo {
	muCache.Lock()
	defer muCache.Unlock()
	nextConnectionID++
	conn := &ConnectionInfo{
		ConnectionString: connString,
		ConnectionID:     nextConnectionID,
		CreatedAt:        time.Now(),
		LastAccessAt:     time.Now(),
		Client:           client,
	}
	connections[conn.ConnectionID] = conn
	return conn
}

func connectionDisconnect(connID int) error {
	muCache.Lock()
	defer muCache.Unlock()
	connInfo := connections[connID]
	if connInfo == nil {
		return fmt.Errorf("disconnect: unknown connection id '%d'", connID)
	}
	delete(connections, connID)
	return connInfo.Client.db.Close()
}

func getConnectionInfoByID(connID int) *ConnectionInfo {
	muCache.Lock()
	defer muCache.Unlock()
	return connections[connID]
}

// TODO: temporary, returns first available connection id
// needs to support multiple connections
// Retruns -1 if there are no connections
func getFirstConnectionId() int {
	muCache.Lock()
	muCache.Unlock()
	for connId, _ := range connections {
		return connId
	}
	return -1
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
	connections = make(map[int]*ConnectionInfo)
}

func isErrNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func execMust(db *sql.DB, q string, args ...interface{}) {
	LogVerbosef("db.Exec(): %s\n", q)
	_, err := db.Exec(q, args...)
	fatalIfErr(err, fmt.Sprintf("db.Exec('%s')", q))
}
