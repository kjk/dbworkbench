package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// A log file that supports rotation based on size threshold

var (
	logToStdout = false
)

// LogFile describes a log file
type LogFile struct {
	sync.Mutex
	path      string
	file      *os.File
	w         io.Writer
	useLogger bool // setting this will log like log package: with date and time
	logger    *log.Logger
}

func (l *LogFile) dirAndBaseName() (string, string) {
	dir := filepath.Dir(l.path)
	base := filepath.Base(l.path)
	return dir, base
}

func (l *LogFile) close() {
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *LogFile) open() (err error) {
	flag := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	l.file, err = os.OpenFile(l.path, flag, 0644)
	if err != nil {
		return err
	}
	_, err = l.file.Stat()
	if err != nil {
		l.file.Close()
		return err
	}
	l.w = l.file
	if l.useLogger {
		l.logger = log.New(l.w, "", log.LstdFlags)
	}
	return err
}

// NewLogFile creates LogRotate
func NewLogFile(path string) (*LogFile, error) {
	res := &LogFile{
		path:      path,
		useLogger: true,
	}
	if err := res.open(); err != nil {
		return nil, err
	}
	return res, nil
}

// Close closes a log file
func (l *LogFile) Close() {
	if l == nil {
		return
	}
	l.Lock()
	l.close()
	l.Unlock()
}

// Print writes to log file
func (l *LogFile) Print(s string) {
	if logToStdout {
		fmt.Print(s)
	}
	if l == nil {
		return
	}
	if l.logger != nil {
		l.logger.Print(s)
	} else {
		l.w.Write([]byte(s))
	}
}

// Printf writes to log file
func (l *LogFile) Printf(format string, arg ...interface{}) {
	l.Print(fmt.Sprintf(format, arg...))
}
