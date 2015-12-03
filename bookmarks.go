package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var bookmarkMutex = &sync.Mutex{}

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
			LogErrorf("Bookmark file creation %v", err)
			return err
		}
		file.Close()
	}

	return nil
}

func readAllBookmarks() (map[string]Bookmark, error) {
	res := map[string]Bookmark{}

	fileData, err := ioutil.ReadFile(bookmarksFilePath())
	if err != nil {
		LogErrorf("Bookmark file read %v", err)
		return res, err
	}

	err = json.Unmarshal(fileData, &res)
	if err != nil {
		LogErrorf("Bookmark unmarshall %v", err)
		return res, err
	}

	return res, nil
}

func readBookmarks() (map[string]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	res, err := readAllBookmarks()
	if err != nil {
		if err == os.ErrNotExist {
			return res, nil
		}
		return res, err
	}
	return res, nil
}

func addBookmark(bookmark Bookmark) (map[string]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	res := map[string]Bookmark{}

	// Path exist
	res, err := readAllBookmarks()
	if err != nil {
		if err == os.ErrNotExist {
			err := ifNotCreateBookmarksFile()
			if err != nil {
				return res, err
			}
		}
		// If the file is empty this will send an error so ignore
		LogInfof("Bookmark file is empty %v", err)
	}

	res[bookmark.Database] = bookmark

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		LogErrorf("Bookmark MarshalIndent %v", err)
		return res, err
	}

	ioutil.WriteFile(bookmarksFilePath(), b, 0644)
	return res, nil
}

func removeBookmark(databaseName string) (map[string]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	res := map[string]Bookmark{}

	// Path exist
	res, err := readAllBookmarks()
	if err != nil {
		if err == os.ErrNotExist {
			return res, nil
		}
		LogErrorf("Bookmark readall %v", err)
		return res, err
	}

	delete(res, databaseName)
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		LogErrorf("Bookmark MarshalIndent %v", err)
		return res, err
	}

	ioutil.WriteFile(bookmarksFilePath(), b, 0644)
	return res, nil
}

type ByDatabaseName []Bookmark

func (s ByDatabaseName) Len() int {
	return len(s)
}
func (s ByDatabaseName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDatabaseName) Less(i, j int) bool {
	return s[i].Database < s[j].Database
}

func sortBookmarks(bookmarks map[string]Bookmark) []Bookmark {
	bookmarkArr := make(ByDatabaseName, 0, len(bookmarks))
	for _, value := range bookmarks {
		bookmarkArr = append(bookmarkArr, value)
	}

	sort.Sort(ByDatabaseName(bookmarkArr))
	return bookmarkArr
}
