package main

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Timing records how long something took to execute
type Timing struct {
	Duration time.Duration
	What     string
}

var (
	timings     []Timing
	inFatal     bool
	logFile     *os.File
	logFileName string // set logFileName to enable loggin
)

func logToFile(s string) {
	if logFileName == "" {
		return
	}

	if logFile == nil {
		var err error
		logFile, err = os.Create(logFileName)
		if err != nil {
			fmt.Printf("logToFile: os.Create('%s') failed with %s\n", logFileName, err)
			os.Exit(1)
		}
	}
	logFile.WriteString(s)
}

func closeLogFile() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// Note: it can say is 32bit on 64bit machine (if 32bit toolset is installed),
// but it'll never say it's 64bit if it's 32bit
func isOS64Bit() bool {
	return runtime.GOARCH == "amd64"
}

func appendTiming(dur time.Duration, what string) {
	t := Timing{
		Duration: dur,
		What:     what,
	}
	timings = append(timings, t)
}

func printTimings() {
	for _, t := range timings {
		fmt.Printf("%s\n    %s\n", t.Duration, t.What)
		logToFile(fmt.Sprintf("%s\n    %s\n", t.Duration, t.What))
	}
}

func printStack() {
	buf := make([]byte, 1024*164)
	n := runtime.Stack(buf, false)
	fmt.Printf("%s", buf[:n])
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	printStack()
	finalizeThings(true)
	os.Exit(1)
}

func fatalif(cond bool, format string, args ...interface{}) {
	if cond {
		if inFatal {
			os.Exit(1)
		}
		inFatal = true
		fmt.Printf(format, args...)
		printStack()
		finalizeThings(true)
		os.Exit(1)
	}
}

func fataliferr(err error) {
	if err != nil {
		fatalf("%s\n", err.Error())
	}
}

func pj(elem ...string) string {
	return filepath.Join(elem...)
}

func replaceExt(path string, newExt string) string {
	ext := filepath.Ext(path)
	return path[0:len(path)-len(ext)] + newExt
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

func toTrimmedLines(d []byte) []string {
	lines := strings.Split(string(d), "\n")
	i := 0
	for _, l := range lines {
		l = strings.TrimSpace(l)
		// remove empty lines
		if len(l) > 0 {
			lines[i] = l
			i++
		}
	}
	return lines[:i]
}

func fileSizeMust(path string) int64 {
	fi, err := os.Stat(path)
	fataliferr(err)
	return fi.Size()
}

func fileCopyMust(dst, src string) {
	in, err := os.Open(src)
	fataliferr(err)
	defer in.Close()

	out, err := os.Create(dst)
	fataliferr(err)

	_, err = io.Copy(out, in)
	cerr := out.Close()
	fataliferr(err)
	fataliferr(cerr)
}

func isNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isGitClean() bool {
	out, err := runExe("git", "status", "--porcelain")
	fataliferr(err)
	s := strings.TrimSpace(string(out))
	return len(s) == 0
}

func removeDirMust(dir string) {
	err := os.RemoveAll(dir)
	fataliferr(err)
}

func removeFileMust(path string) {
	if !fileExists(path) {
		return
	}
	err := os.Remove(path)
	fataliferr(err)
}

// Version must be in format x.y.z
func verifyCorrectVersionMust(ver string) {
	parts := strings.Split(ver, ".")
	fatalif(len(parts) == 0 || len(parts) > 3, "%s is not a valid version number", ver)
	for _, part := range parts {
		fatalif(!isNum(part), "%s is not a valid version number", ver)
	}
}

func getGitSha1Must() string {
	out, err := runExe("git", "rev-parse", "HEAD")
	fataliferr(err)
	s := strings.TrimSpace(string(out))
	fatalif(len(s) != 40, "getGitSha1Must(): %s doesn't look like sha1\n", s)
	return s
}

func dataSha1Hex(d []byte) string {
	sha1 := sha1.Sum(d)
	return fmt.Sprintf("%x", sha1[:])
}

func fileSha1Hex(path string) (string, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	sha1 := sha1.Sum(d)
	return fmt.Sprintf("%x", sha1[:]), nil
}

func httpDlMust(uri string) []byte {
	res, err := http.Get(uri)
	fataliferr(err)
	d, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fataliferr(err)
	return d
}

func httpDlToFileMust(uri string, path string, sha1Hex string) {
	if fileExists(path) {
		sha1File, err := fileSha1Hex(path)
		fataliferr(err)
		sha1Hex = strings.ToUpper(sha1Hex)
		sha1File = strings.ToUpper(sha1File)
		fatalif(sha1File != sha1Hex, "file '%s' exists but has sha1 of %s and we expected %s", path, sha1File, sha1Hex)
		return
	}
	fmt.Printf("Downloading '%s'\n", uri)
	d := httpDlMust(uri)
	sha1File := dataSha1Hex(d)
	fatalif(sha1File != sha1Hex, "downloaded '%s' but it has sha1 of %s and we expected %s", uri, sha1File, sha1Hex)
	err := ioutil.WriteFile(path, d, 0755)
	fataliferr(err)
}

// ZipDirectory creates a zip file out of directory
func ZipDirectory(dirToZip, zipPath string) error {
	stat, err := os.Stat(dirToZip)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dirToZip)
	}

	baseDir := filepath.Base(dirToZip)

	zipfile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(dirToZip, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate

		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, dirToZip))

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		return err
	})

	return err
}

func zipFileName(path, baseDir string) string {
	fatalif(!strings.HasPrefix(path, baseDir), "'%s' doesn't start with '%s'", path, baseDir)
	n := len(baseDir)
	path = path[n:]
	if path[0] == '/' || path[0] == '\\' {
		path = path[1:]
	}
	// always use unix path separator inside zip files because that's what
	// the browser uses in url and we must match that
	return strings.Replace(path, "\\", "/", -1)
}

func addZipFileMust(zw *zip.Writer, path, zipName string) {
	fi, err := os.Stat(path)
	fataliferr(err)
	fih, err := zip.FileInfoHeader(fi)
	fataliferr(err)
	fih.Name = zipName
	fih.Method = zip.Deflate
	d, err := ioutil.ReadFile(path)
	fataliferr(err)
	fw, err := zw.CreateHeader(fih)
	fataliferr(err)
	_, err = fw.Write(d)
	fataliferr(err)
	// fw is just a io.Writer so we can't Close() it. It's not necessary as
	// it's implicitly closed by the next Create(), CreateHeader()
	// or Close() call on zip.Writer
}

func addZipDirMust(zw *zip.Writer, dir, baseDir string) {
	dirsToVisit := []string{dir}
	for len(dirsToVisit) > 0 {
		dir = dirsToVisit[0]
		dirsToVisit = dirsToVisit[1:]
		files, err := ioutil.ReadDir(dir)
		fataliferr(err)
		for _, fi := range files {
			name := fi.Name()
			path := filepath.Join(dir, name)
			if fi.IsDir() {
				dirsToVisit = append(dirsToVisit, path)
			} else if fi.Mode().IsRegular() {
				zipName := zipFileName(path, baseDir)
				if !isBlacklisted(zipName) {
					addZipFileMust(zw, path, zipName)
				}
			}
		}
	}
}
