package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/kjk/u"
)

var bookmarkMutex = sync.Mutex{}

// Bookmark defines info about a database connection
type Bookmark struct {
	// DbType string `json:"dbtype"` // postgres, mysql etc.
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Ssl      string `json:"ssl"`
}

func bookmarksFilePath() string {
	return filepath.Join(getDataDir(), "bookmarks.json")
}

// Unnamed Connection is a predefined keyword which is used in frontend
// In frontend the Bookmark.Database is the bookmark name which needs to
// be set
func createBookmarkForUnnamedConnection() {
	if u.PathExists(bookmarksFilePath()) {
		return
	}
	addBookmark(Bookmark{Database: "Unnamed Connection"})
}

func readAllBookmarks() (map[string]Bookmark, error) {
	res := map[string]Bookmark{}

	fileData, err := ioutil.ReadFile(bookmarksFilePath())
	if err != nil {
		if err == os.ErrNotExist {
			return res, nil
		}
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

	bookmarks, err := readAllBookmarks()
	if err != nil {
		return nil, err
	}
	for _, b := range bookmarks {
		pwd := b.Password
		b.Password, err = decrypt(pwd)
		if err != nil {
			LogErrorf("load1s('%s') failed with '%s'\n", pwd, err)
			b.Password = ""
		}
	}
	return bookmarks, nil
}

func addBookmark(bookmark Bookmark) (map[string]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmark.Password = encrypt(bookmark.Password)

	bookmarks, err := readAllBookmarks()
	if err != nil {
		// If the file is empty this will send an error so ignore
		LogInfof("Bookmark file is empty %v", err)
	}

	bookmarks[bookmark.Database] = bookmark

	b, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		LogErrorf("Bookmark MarshalIndent %v", err)
		return bookmarks, err
	}

	err = ioutil.WriteFile(bookmarksFilePath(), b, 0644)
	if err != nil {
		LogErrorf("ioutil.WriteFile() failed with '%s'\n", err)
	}
	return bookmarks, nil
}

func removeBookmark(databaseName string) (map[string]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmarks, err := readAllBookmarks()
	if err != nil {
		LogErrorf("readAllBookmarks() failed with '%s'\n", err)
		return bookmarks, err
	}

	delete(bookmarks, databaseName)
	b, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		LogErrorf("Bookmark MarshalIndent %v", err)
		return bookmarks, err
	}

	ioutil.WriteFile(bookmarksFilePath(), b, 0644)
	return bookmarks, nil
}

// ByDatabaseName is for sorting by database name
type ByDatabaseName []Bookmark

func (s ByDatabaseName) Len() int           { return len(s) }
func (s ByDatabaseName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByDatabaseName) Less(i, j int) bool { return s[i].Database < s[j].Database }

func sortBookmarks(bookmarks map[string]Bookmark) []Bookmark {
	bookmarkArr := make(ByDatabaseName, 0, len(bookmarks))
	for _, value := range bookmarks {
		bookmarkArr = append(bookmarkArr, value)
	}

	sort.Sort(ByDatabaseName(bookmarkArr))
	return bookmarkArr
}
