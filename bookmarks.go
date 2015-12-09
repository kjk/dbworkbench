package main

import (
	"encoding/json"
	"fmt"
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
	ID       int    `json:"id"`
	Type     string `json:"type"` // postgres, mysql etc.
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
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
	b := Bookmark{
		ID:       1,
		Type:     "postgres",
		Database: "Unnamed Connection",
	}
	addBookmark(b)
}

// find the lowest available id, starting with 1
func findBookmarkID(bookmarks []Bookmark) int {
	n := len(bookmarks)
	if n == 0 {
		return 1
	}
	a := make([]int, n, n)
	for i := 0; i < n; i++ {
		a[i] = bookmarks[i].ID
	}
	sort.Ints(a)
	newID := a[0]
	// find a gap in ids
	for i := 1; i < n; i++ {
		newID++
		if a[i] != newID {
			return newID
		}
	}
	// or return the next one
	return newID + 1
}

func readAllBookmarks() ([]Bookmark, error) {
	var res []Bookmark

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

	for _, b := range res {
		pwd := b.Password
		b.Password, err = decrypt(pwd)
		if err != nil {
			LogInfof("decrypted '%s' => '%s'\n", pwd, b.Password)
			LogErrorf("decrypt('%s') failed with '%s'\n", pwd, err)
			b.Password = ""
		}
	}

	return res, nil
}

func readBookmarks() ([]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmarks, err := readAllBookmarks()
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// checks if database type is correct. Returns normalized
// name of "" if not valid
func validateDatabaseType(s string) string {
	if s == "" {
		return "postgres"
	}
	if s == "postgres" || s == "pg" {
		return "postgres"
	}
	if s == "mysql" {
		return s
	}
	return ""
}

func addBookmark(bookmark Bookmark) ([]Bookmark, error) {
	tp := validateDatabaseType(bookmark.Type)
	if tp == "" {
		return nil, fmt.Errorf("invalid bookmark type '%s'", bookmark.Type)
	}
	bookmark.Type = tp
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmark.Password = encrypt(bookmark.Password)

	bookmarks, err := readAllBookmarks()
	if err != nil {
		// If the file is empty this will send an error so ignore
		LogInfof("Bookmark file is empty %v", err)
	}

	bookmarks = append(bookmarks, bookmark)
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

func removeBookmarkByID(arr []Bookmark, id int) []Bookmark {
	n := len(arr)
	for i := 0; i < n; i++ {
		if arr[i].ID == id {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func removeBookmark(id int) ([]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmarks, err := readAllBookmarks()
	if err != nil {
		LogErrorf("readAllBookmarks() failed with '%s'\n", err)
		return bookmarks, err
	}

	bookmarks = removeBookmarkByID(bookmarks, id)
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

func sortBookmarks(bookmarks []Bookmark) []Bookmark {
	sort.Sort(ByDatabaseName(bookmarks))
	return bookmarks
}
