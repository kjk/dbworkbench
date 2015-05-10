package main

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

/*
TODO: implement
*/

type Bookmark struct {
	Url      string `json:"url"`      // Postgres connection URL
	Host     string `json:"host"`     // Server hostname
	Port     string `json:"port"`     // Server port
	User     string `json:"user"`     // Database user
	Password string `json:"password"` // User password
	Database string `json:"database"` // Database name
	Ssl      string `json:"ssl"`      // Connection SSL mode
}

func bookmarksPath() string {
	path, _ := homedir.Dir()
	return fmt.Sprintf("%s/.pgweb/bookmarks", path)
}

func readAllBookmarks(path string) (map[string]Bookmark, error) {
	res := map[string]Bookmark{}
	return res, nil

}
