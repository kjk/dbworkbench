package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	blaclisted = []string{
		"s/dist/bundle.js.map",
		".gitkeep",
		// TODO: eventually simply delete files we don't use
		// we use bootstrap-flatly.css
		"s/css/bootstrap.css",
		"s/css/bootstrap-lumen.css",
		"s/css/bootstrap-paper.css",
		"s/css/bootstrap-simplex.css",
		"s/css/normalize.css",
		"s/css/primer.css",
		"s/css/pure-min.css",
	}
)

func isBlacklisted(path string) bool {
	for _, s := range blaclisted {
		if path == s {
			return true
		}
	}
	return false
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

func createResourcesZip() {
	path := "dbworkbench.dat"
	f, err := os.Create(path)
	fataliferr(err)
	defer f.Close()
	zw := zip.NewWriter(f)
	currDir, err := os.Getwd()
	fataliferr(err)
	dir := filepath.Join(currDir, "s")
	addZipDirMust(zw, dir, currDir)
	err = zw.Close()
	fataliferr(err)
}
