package main

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
)

var (
	usageFile                  *os.File
	usageFileMutex             sync.Mutex
	usageFileWriteFailedBefore bool
)

func openUsageFileMust() {
	var err error
	path := filepath.Join(getDataDir(), "usage.txt")
	usageFile, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		LogFatalf("os.OpenFile('%s') failed with '%s'\n", path, err)
	}
}

func sanitizeUsage(d []byte) []byte {
	// normalize newlines
	d = bytes.Replace(d, []byte{'\r', '\n'}, []byte{'\n'}, -1)
	d = bytes.Replace(d, []byte{'\r'}, []byte{'\n'}, -1)
	// remove empty-lines ('\n\n') until none are left
	prevSize := len(d) - 1 // just to make it different
	for prevSize != len(d) {
		prevSize = len(d)
		d = bytes.Replace(d, []byte{'\n', '\n'}, []byte{'\n'}, -1)
	}
	return d
}

func recordUsage(d []byte) {
	// since we separate records with empty lines, make sure
	// that data itself doesn't contain empty lines
	d = sanitizeUsage(d)
	usageFileMutex.Lock()
	defer usageFileMutex.Unlock()
	_, err := usageFile.Write(d)
	// prevent flooding the log in case this starts happening often
	if err != nil && !usageFileWriteFailedBefore {
		usageFileWriteFailedBefore = true
		LogErrorf("usageFile.Write() failed with '%s'\n", err)
	}
}
