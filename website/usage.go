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
	LogInfof("usge file: '%s'\n", path)
}

func sanitizeUsage(d []byte) []byte {
	// normalize newlines
	d = bytes.Replace(d, []byte{'\r', '\n'}, []byte{'\n'}, -1)
	d = bytes.Replace(d, []byte{'\r'}, []byte{'\n'}, -1)
	// remove empty-lines ('\n\n') until none are left
	prevSize := -1 // at first iteration must be != len(d)
	for prevSize != len(d) {
		prevSize = len(d)
		d = bytes.Replace(d, []byte{'\n', '\n'}, []byte{'\n'}, -1)
	}
	n := len(d)
	if n == 0 {
		return d
	}
	last := d[n-1]
	d = append(d, '\n')
	if last != '\n' {
		d = append(d, '\n')
	}
	return d
}

func recordUsage(d []byte) {
	// since we separate records with empty lines, make sure
	// that data itself doesn't contain empty lines
	d = sanitizeUsage(d)
	if len(d) == 0 {
		return
	}
	//LogInfof("recordUsage: len(d) = %d\n", len(d))
	usageFileMutex.Lock()
	defer usageFileMutex.Unlock()
	_, err := usageFile.Write(d)
	// prevent flooding the log in case this starts happening often
	if err != nil && !usageFileWriteFailedBefore {
		usageFileWriteFailedBefore = true
		LogErrorf("usageFile.Write() failed with '%s'\n", err)
	}
}
