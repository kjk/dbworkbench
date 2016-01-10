package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const hdr = `// +build embeded_resources

package main

var resourcesZipData = []byte{
`

var (
	blacklisted = []string{
		filepath.Join("s", "dist", "bundle.js"),
		filepath.Join("s", "dist", "main2.css"),
		filepath.Join("s", "fake_data_1.js"),
		".gitkeep",
	}
)

func isBlacklisted(path string) bool {
	if strings.HasSuffix(path, ".map") {
		return true
	}
	for _, s := range blacklisted {
		if path == s {
			return true
		}
	}
	return false
}

func createResourcesZip(path string) {
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

func genHexLine(f *os.File, d []byte, off, n int) {
	f.WriteString("\t")
	for i := 0; i < n; i++ {
		b := d[off+i]
		fmt.Fprintf(f, "0x%02x,", b)
	}
	f.WriteString("\n")
}

func genResourcesGo(goPath, dataPath string) {
	d, err := ioutil.ReadFile(dataPath)
	fataliferr(err)
	f, err := os.Create(goPath)
	fataliferr(err)
	defer f.Close()
	f.WriteString(hdr)

	nPerLine := 16
	nFullLines := len(d) / nPerLine
	nLastLine := len(d) % nPerLine
	n := 0
	for i := 0; i < nFullLines; i++ {
		genHexLine(f, d, n, nPerLine)
		n += nPerLine
	}
	genHexLine(f, d, n, nLastLine)
	f.WriteString("}\n")
}

func genResources() {
	zipPath := "dbherohelper.zip"
	createResourcesZip(zipPath)
	goPath := "resources.go"
	genResourcesGo(goPath, zipPath)
}
