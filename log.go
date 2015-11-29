package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	logFile *LogFile

	dot               = []byte(".")
	centerDot         = []byte("·")
	logVerbosityLevel int32

	logToStdout = false
)

// IncLogVerbosity increases log verbosity
// the idea of verbose logging is to provide a way to turn detailed logging
// on a per-request basis. This is an approximate solution: since there is
// no per-gorutine context, we use a shared variable that is increased at request
// beginning and decreased at end. We might get additional logging from other
// gorutines. It's much simpler than an alternative, like passing a logger
// to every function that needs to log
func IncLogVerbosity() {
	atomic.AddInt32(&logVerbosityLevel, 1)
}

// DecLogVerbosity decreases log verbosity
func DecLogVerbosity() {
	atomic.AddInt32(&logVerbosityLevel, -1)
}

// IsVerboseLogging returns true if verbose logging is turned on
func IsVerboseLogging() bool {
	return atomic.LoadInt32(&logVerbosityLevel) > 0
}

// StartVerboseLog is meant to be used like this:
// if StartVerboseLog(r.URL) {
//      defer DecLogVerbosity()
// }
func StartVerboseLog(u *url.URL) bool {
	// "vl" stands for "verbose logging" and any value other than empty string
	// truns it on
	if u.Query().Get("vl") != "" {
		IncLogVerbosity()
		return true
	}
	return false
}

// StopVerboseLog is for name parity with StartVerboseLog()
func StopVerboseLog() {
	DecLogVerbosity()
}

// OpenLogMust opens a log file
func OpenLogMust(fileName string) *LogFile {
	path := filepath.Join(getLogDir(), fileName)
	fmt.Printf("log file: %s\n", path)
	logTmp, err := NewLogFile(path)
	if logTmp == nil {
		log.Fatalf("OpenLogMust: NewLogFile(%q) failed with %q\n", path, err)
	}
	return logTmp
}

func dumpStrArr(arr []string) {
	for _, s := range arr {
		fmt.Printf("%s\n", s)
	}
}

// RemoveOldLogFiles removes old log files
func RemoveOldLogFiles() {
	files, err := ioutil.ReadDir(getLogDir())
	if err != nil {
		return
	}
	var fileNames []string
	for _, fi := range files {
		name := fi.Name()
		if strings.HasPrefix(name, "log-") && strings.HasSuffix(name, ".txt") {
			fileNames = append(fileNames, name)
		}
	}
	sort.Strings(fileNames)
	// keep last 10 log files for easier debugging
	maxToKeep := 10
	n := len(fileNames)
	if n <= maxToKeep {
		return
	}
	fileNames = fileNames[n-maxToKeep:]
	for _, name := range fileNames {
		path := filepath.Join(getLogDir(), name)
		os.Remove(path)
	}
}

// OpenLogFiles open error and info log files
func OpenLogFiles() {
	t := time.Now()
	fileName := t.Format("log-06-01-02_15_04_05.txt")
	logFile = OpenLogMust(fileName)
}

// CloseLogFiles closes log files
func CloseLogFiles() {
	logFile.Close()
	logFile = nil
}

func functionFromPc(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//      runtime/debug.*T·ptrmethod
	// and want
	//      *T.ptrmethod
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return string(name)
}

// LogFatalf is like log.Fatalf() but also pre-pends name of the caller,
// so that we don't have to do that manually in every log statement
func LogFatalf(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	fmt.Print(s)
	logFile.Print(s)

	CloseLogFiles()
	log.Fatal(s)
}

// LogErrorf is for logging things that are unexpected but not fatal
// Automatically pre-pends name of the function calling the log function
func LogErrorf(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}

	logFile.Print(s)
}

// LogError logs an error
func LogError(s string) {
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	logFile.Print(s)
}

// LogInfof is for logging of misc non-error things
func LogInfof(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	logFile.Print(s)
}

// LogVerbosef is meant for detailed information that is only enabled on
// a per request basis
func LogVerbosef(format string, arg ...interface{}) {
	if !IsVerboseLogging() {
		return
	}
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	logFile.Print(s)
}

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
