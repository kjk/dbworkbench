package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kjk/dbworkbench/ga_event"
)

func get(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject non-GET methods
	return f
}

func post(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject non-POST methods
	return f
}

func hasUser(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject if no user
	return f
}

func hasConnection(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject if no connection
	return f
}
func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: write me
	/*
	  data, err := asset(path)

	  if err != nil {
	    c.String(400, err.Error())
	    return
	  }

	  c.Data(200, "text/html; charset=utf-8", data)
	*/
}

// GET /
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri != "/" {
		http.NotFound(w, r)
		return
	}
	serveStatic(w, r, "s/index_react.html")
}

// GET /s/:path
func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/s/"):]
	serveStatic(w, r, path)
}

// POST /api/connect
func handleConnect(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/history
func handleHistory(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/bookmarks
func handleBookmarks(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/databases
func handleGetDatabases(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/connection
func handleConnectionInfo(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/activity
func handleActivity(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/schemas
func handleGetSchemas(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/tables
func handleGetTables(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

// GET /api/tables/:table/:action
func handleTablesDispatch(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
}

func registerHTTPHandlers() {
	http.HandleFunc("/", get(handleIndex))
	http.HandleFunc("/s/", get(handleStatic))
	http.HandleFunc("/api/connect", post(hasUser(handleConnect)))
	http.HandleFunc("/api/history", get(hasUser(handleHistory)))
	http.HandleFunc("/api/bookmarks", get(hasUser(handleBookmarks)))

	http.HandleFunc("/api/databases", get(hasUser(hasConnection(handleGetDatabases))))
	http.HandleFunc("/api/connection", get(hasUser(hasConnection(handleConnectionInfo))))
	http.HandleFunc("/api/activity", get(hasUser(hasConnection(handleActivity))))
	http.HandleFunc("/api/schemas", get(hasUser(hasConnection(handleGetSchemas))))
	http.HandleFunc("/api/tables", get(hasUser(hasConnection(handleGetTables))))
	http.HandleFunc("/api/tables/", get(hasUser(hasConnection(handleTablesDispatch))))
}

func setupRoutes(router *gin.Engine) {

	api := router.Group("/api")
	{
		api.Use(ApiMiddleware())

		api.GET("/tables/:table", API_GetTable)
		api.GET("/tables/:table/rows", API_GetTableRows)
		api.GET("/tables/:table/info", API_GetTableInfo)
		api.GET("/tables/:table/indexes", API_TableIndexes)

		api.GET("/query", API_RunQuery)
		api.POST("/query", API_RunQuery)
		api.GET("/explain", API_ExplainQuery)
		api.POST("/explain", API_ExplainQuery)
	}
}

func startWebServer() {
	registerHTTPHandlers()
	httpAddr := fmt.Sprintf(":%v", options.HttpPort)
	fmt.Printf("Started running on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListendAndServer() failed with %s\n", err)
	}
	fmt.Printf("Exited\n")
}

func startGinServer() {
	router := gin.Default()
	router.Use(ga_event.GALogger("UA-62336732-1", "databaseworkbench.com"))

	// Enable HTTP basic authentication only if both user and password are set
	if options.AuthUser != "" && options.AuthPass != "" {
		auth := map[string]string{options.AuthUser: options.AuthPass}
		router.Use(gin.BasicAuth(auth))
	}

	setupRoutes(router)

	fmt.Println("Starting server...")
	go func() {
		err := router.Run(fmt.Sprintf("%v:%v", options.HttpHost, options.HttpPort))
		if err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
