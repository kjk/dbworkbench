package main

import (
 	"encoding/json"
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
)

// TODO encript password

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

func ifNotCreateBookmarksFile() error {
    if _, err := os.Stat(bookmarksFilePath()); os.IsNotExist(err) {
        // Path does not exist
        file, err := os.Create(bookmarksFilePath()) // Maybe use os.NewFile?
        if err != nil {
            fmt.Println(err)
            return err
        }
        defer file.Close()
    }

    return nil
}

func readAllBookmarks() (map[string]Bookmark, error) {
    ifNotCreateBookmarksFile()

    res := map[string]Bookmark{}

	fileData, err := ioutil.ReadFile(bookmarksFilePath())
    if err != nil {
        fmt.Printf("File error: %v\n", err)
        return res, err
    }

    err = json.Unmarshal(fileData, &res)
    if err != nil {
        fmt.Printf("JSON Parse error: %v\n", err)
    	return res, err
    }

	return res, nil
}

func addBookmark(bookmark Bookmark) (map[string]Bookmark, error) {
    ifNotCreateBookmarksFile()

    res := map[string]Bookmark{}


    // Path exist
    res, err := readAllBookmarks()
    if err != nil {
        // If the file is empty this will send an error so ignore
        fmt.Printf("readAllBookmarks error: %v\n", err)
        // return res, err
    }

    res[bookmark.Database] = bookmark

    b, err := json.MarshalIndent(res, "", "  ")
    if err != nil {
     	return res, err
 	}

    ioutil.WriteFile(bookmarksFilePath(), b, 0644)
    return res, nil
}

func removeBookmark(databaseName string) (map[string]Bookmark, error) {
    ifNotCreateBookmarksFile()

    res := map[string]Bookmark{}

    // Path exist
    res, err := readAllBookmarks()
    if err != nil {
        return res, err
    }

    delete(res, databaseName)
    b, err := json.MarshalIndent(res, "", "  ")
    if err != nil {
    	return res, err
    }

    ioutil.WriteFile(bookmarksFilePath(), b, 0644)
    return res, nil
}
