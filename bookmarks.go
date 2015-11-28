package main

import "path/filepath"

/*
TODO: implement.
*/

// Bookmark defines info about a database
type Bookmark struct {
	URL      string `json:"url"`      // Postgres connection URL
	Host     string `json:"host"`     // Server hostname
	Port     string `json:"port"`     // Server port
	User     string `json:"user"`     // Database user
	Password string `json:"password"` // User password
	Database string `json:"database"` // Database name
	Ssl      string `json:"ssl"`      // Connection SSL mode
}

func bookmarksFilePath() string {
	return filepath.Join(getDataDir(), "bookmarks.json")
}

func readAllBookmarks() (map[string]Bookmark, error) {
	res := map[string]Bookmark{}
	return res, nil
}
