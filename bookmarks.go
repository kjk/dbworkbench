package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
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

// find the lowest available id, starting with 1
func generateNewBookmarkID(bookmarks []Bookmark) int {
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

// must be called while holding bookmarkMutex lock
func readBookmarksDecryptPwd() ([]Bookmark, error) {
	var res []Bookmark

	path := bookmarksFilePath()
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return res, nil
		}
		LogErrorf("ioutil.ReadFile('%s') failed with '%s'\n", path, err)
		return nil, err
	}

	err = json.Unmarshal(fileData, &res)
	if err != nil {
		LogErrorf("json.Unmarshall() failed with '%s'\n", err)
		return nil, err
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

// must be called while holding bookmarkMutex lock
// encrypts passwords and writes bookmarks as json to bookmarks.json
// returns original bookmarks + error for convenience of the caller
func writeBookmarksEncryptPwd(bookmarks []Bookmark) ([]Bookmark, error) {
	// encrypt passwords in a copy
	n := len(bookmarks)
	a := make([]Bookmark, n, n)
	for i := 0; i < n; i++ {
		a[i] = bookmarks[i]
		a[i].Password = encrypt(a[i].Password)
	}

	d, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		LogErrorf("json.MarshalIndent() failed with '%s'", err)
		return bookmarks, err
	}

	path := bookmarksFilePath()
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		LogErrorf("ioutil.WriteFile('%s') failed with '%s'\n", path, err)
	}
	return bookmarks, err
}

func readBookmarks() ([]Bookmark, error) {
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	return readBookmarksDecryptPwd()
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

func findBookmarkIndexByID(arr []Bookmark, id int) int {
	n := len(arr)
	for i := 0; i < n; i++ {
		if arr[i].ID == id {
			return i
		}
	}
	return -1
}

func addBookmark(bookmark Bookmark) ([]Bookmark, error) {
	tp := validateDatabaseType(bookmark.Type)
	if tp == "" {
		return nil, fmt.Errorf("invalid bookmark type '%s'", bookmark.Type)
	}
	bookmark.Type = tp
	bookmarkMutex.Lock()
	defer bookmarkMutex.Unlock()

	bookmarks, err := readBookmarksDecryptPwd()
	if err != nil {
		return nil, err
	}

	idx := findBookmarkIndexByID(bookmarks, bookmark.ID)

	if idx == -1 {
		bookmark.ID = generateNewBookmarkID(bookmarks)
		bookmarks = append(bookmarks, bookmark)
	} else {
		bookmarks[idx] = bookmark
	}
	return writeBookmarksEncryptPwd(bookmarks)
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

	bookmarks, err := readBookmarksDecryptPwd()
	if err != nil {
		return nil, err
	}

	bookmarks = removeBookmarkByID(bookmarks, id)
	return writeBookmarksEncryptPwd(bookmarks)
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
