package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func testMysqlInfo() {
	// ${user}:${pwd}@tcp(${ip}:${port})/{dbName}?parseTime=true
	uri := os.Getenv("MYSQL_URL")
	if uri == "" {
		log.Fatal("MYSQL_URL env variable is not defined. in bash, do:\nMYSQL_URL=${uri} ./scripts.run.sh\n")
	}
	c, err := NewClientMysqlFromURL(uri)
	if err != nil {
		fmt.Printf("NewClientMysqlFromURL('%s') failed with '%s'\n", uri, err)
		return
	}
	db := c.Connection()
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Printf("c.Test() failed with '%s', uri: '%s'\n", err, uri)
		return
	}

	i, err := c.Info()
	if err != nil {
		fmt.Printf("c.Info() failed with '%s'\n", err)
		return
	}
	i.Dump()

	dbs, err := c.Databases()
	if err != nil {
		fmt.Printf("c.Databases() failed with '%s'\n", err)
		return
	}
	fmt.Printf("%d datbases:\n", len(dbs))
	for _, db := range dbs {
		fmt.Printf("  '%s'\n", db)
	}

	tables, err := c.Tables()
	if err != nil {
		fmt.Printf("c.Tables() failed with '%s'\n", err)
		return
	}
	fmt.Printf("%d tables:\n", len(tables))
	for _, s := range tables {
		fmt.Printf("  '%s'\n", s)
	}

	/*
		for _, t := range tables {
			schema, err := c.Table(t)
			if err != nil {
				fmt.Printf("c.TableRows('%s') failed with '%s'\n", t, err)
				return
			}
			fmt.Printf("Table: '%s'\n", t)
			schema.DumpFull()
		}
	*/
	//dumpQuery(db, mysqlActivityStmt)
	//dumpQuery(db, `SELECT VARIABLE_NAME, VARIABLE_VALUE FROM INFORMATION_SCHEMA.GLOBAL_VARIABLES`)
	dumpQuery(db, `SELECT * from information_schema.tables;`)
}

func dumpQuery(db *sqlx.DB, query string) {
	res, err := dbQuery(db, query)
	if err != nil {
		fmt.Printf("dbQuery() failed with '%s'\n", err)
		return
	}
	res.DumpFull()
}

func testMysqlTimeoutWithURI(uri string) {
	timeStart := time.Now()
	fmt.Printf("Starting NewClientMysqlFromURL()...")
	c, err := NewClientMysqlFromURL(uri)
	if err != nil {
		fmt.Printf("\nNewClientMysqlFromURL('%s') failed with '%s' after %s\n", uri, err, time.Since(timeStart))
		return
	}
	db := c.Connection()
	defer db.Close()
	timeStartPing := time.Now()
	fmt.Printf("\nstarting ping...")
	if err = db.Ping(); err != nil {
		fmt.Printf("c.Test() failed with '%s', uri: '%s' after %s\n", err, uri, time.Since(timeStartPing))
	}
	fmt.Printf("\ntook %s\n", time.Since(timeStartPing))
}

func testPgTimeoutWithURI(uri string) {
	timeStart := time.Now()
	fmt.Printf("Starting NewClientPgFromURL()...")
	c, err := NewClientPgFromURL(uri)
	if err != nil {
		fmt.Printf("\nNewClientMysqlFromURL('%s') failed with '%s' after %s\n", uri, err, time.Since(timeStart))
		return
	}

	db := c.Connection()
	defer db.Close()
	timeStartPing := time.Now()
	fmt.Printf("\nstarting ping...")
	if err = db.Ping(); err != nil {
		fmt.Printf("c.Test() failed with '%s', uri: '%s' after %s\n", err, uri, time.Since(timeStartPing))
	}
	fmt.Printf("\ntook %s\n", time.Since(timeStartPing))
}

func testMysqlTimeout() {
	// test for timeout by trying to connect to a wrong port
	testMysqlTimeoutWithURI("root:foo@tcp(173.194.251.111:5432)/hello")
}

func testPostgresTimeout() {
	// test for timeout by trying to connect to a wrong port
	testPgTimeoutWithURI("postgres://root:foo@173.194.251.111:5432/hello")
}
