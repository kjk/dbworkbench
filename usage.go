package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record information about program usage like counting how many times
// a database is opened, how many queries are executed etc.
// When desktop app is started, it reads this file, sends it to the website
// for analytics with auto update request.
// The file is a json object (for flexibility) per one line

// UsageEvent describes a user action we track. We use short json names for compactness
type UsageEvent struct {
	// for compactness, we record time as unix epoch i.e. number of
	// seconds since jan 1 1970 UTC
	Timestamp int64  `json:"t"`
	EventType string `json:"e"`
}

const (
	eventDatabaseOpened = "do"
	eventQueryExecuted  = "qe"
)

var (
	usageFile        *os.File
	usageFileEncoder *json.Encoder
	usageFileMutex   sync.Mutex
)

func openUsageFileMust() {
	var err error
	path := filepath.Join(getDataDir(), "usage.json")
	usageFile, err = os.Create(path)
	if err != nil {
		LogFatalf("os.Create('%s') failed with '%s'\n", path, err)
	}
	usageFileEncoder = json.NewEncoder(usageFile)
}

func closeUsageFile() {
	if usageFile != nil {
		usageFile.Close()
		usageFile = nil
	}
}

func recordEvent(e *UsageEvent) {
	usageFileMutex.Lock()
	defer usageFileMutex.Unlock()
	err := usageFileEncoder.Encode(e)
	if err != nil {
		LogErrorf("usageFileEncoder.Encode() failed with '%s'\n", err)
	}
}

func recordDatabaseOpened() {
	e := UsageEvent{
		Timestamp: time.Now().Unix(),
		EventType: eventDatabaseOpened,
	}
	recordEvent(&e)
}

func recordQueryExecuted() {
	e := UsageEvent{
		Timestamp: time.Now().Unix(),
		EventType: eventQueryExecuted,
	}
	recordEvent(&e)
}
