package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/kjk/u"
)

const (
	s3Dir = "software/databaseworkbench/"
)

var (
	flgNoCleanCheck     bool
	flgUpload           bool
	flgUploadAutoUpdate bool
	flgGenResources     bool

	programVersionWin string
	programVersionMac string
	programVersion    string
	certPath          string
	cachedSecrets     *Secrets
	innoSetupPath     string
)

// Secrets defines secrets
type Secrets struct {
	AwsSecret        string
	AwsAccess        string
	CertPwd          string
	NotifierEmail    string
	NotifierEmailPwd string
}

func finalizeThings(crashed bool) {
	// nothing for us
}

func parseCmdLine() {
	// -no-clean-check is useful when testing changes to this build script
	flag.BoolVar(&flgNoCleanCheck, "no-clean-check", false, "allow running if repo has changes (for testing build script)")
	flag.BoolVar(&flgUpload, "upload", false, "if true, upload the files")
	flag.BoolVar(&flgGenResources, "gen-resources", false, "if true, genereates resources.go file")
	flag.Parse()
}

func isGitCleanMust() bool {
	out, err := runExe("git", "status", "--porcelain")
	fataliferr(err)
	s := strings.TrimSpace(string(out))
	return len(s) == 0
}

func verifyGitCleanMust() {
	if !flgUpload || flgNoCleanCheck {
		return
	}
	fatalif(!isGitCleanMust(), "git has unsaved changes\n")
}

func readSecretsMust(path string) *Secrets {
	if cachedSecrets != nil {
		return cachedSecrets
	}
	d, err := ioutil.ReadFile(path)
	fatalif(err != nil, "readSecretsMust(): error %s reading file '%s'\n", err, path)
	var s Secrets
	err = json.Unmarshal(d, &s)
	fatalif(err != nil, "readSecretsMust(): failed to json-decode file '%s'. err: %s, data:\n%s\n", path, err, string(d))
	cachedSecrets = &s
	return cachedSecrets
}

func verifyHasSecretsMust() {
	certPath = pj("scripts", "cert.pfx")
	if !fileExists(certPath) {
		certPath = pj("..", "..", "..", "..", "..", "sumatrapdf", "scripts", "cert.pfx")
		fatalif(!fileExists(certPath), "didn't find cert.pfx in scripts/ or ../../../sumatrapdf/scripts/")
	}
	absPath, err := filepath.Abs(certPath)
	fataliferr(err)
	certPath = absPath
	fmt.Printf("signing certificate: '%s'\n", certPath)

	secretsPath := pj("scripts", "secrets.json")
	if !fileExists(secretsPath) {
		secretsPath = pj("..", "..", "..", "..", "..", "sumatrapdf", "scripts", "secrets.json")
		fatalif(!fileExists(secretsPath), "didn't find secrets.json in scripts/ or ../../../sumatrapdf/scripts/")
	}
	secrets := readSecretsMust(secretsPath)
	if flgUpload || flgUploadAutoUpdate {
		// when uploading must have s3 credentials
		s3SetSecrets(secrets.AwsAccess, secrets.AwsSecret)
	}
}

func isWin() bool {
	return runtime.GOOS == "windows"
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

// given 1.1.0 returns "1.1" i.e. removes all trailing '0' and '.'
func cleanVer(s string) string {
	return strings.TrimRight(s, "0.")
}

func extractVersionWinMust() {
	path := filepath.Join("win", "DatabaseWorkbench", "Properties", "AssemblyInfo.cs")
	d, err := ioutil.ReadFile(path)
	fataliferr(err)
	r := regexp.MustCompile(`(?mi)AssemblyFileVersion\("([^"]+)`)
	s := string(d)
	res := r.FindStringSubmatch(s)
	fatalif(len(res) != 2, "didn't find AssemblyFileVersion in:\n'%s'\n", s)
	programVersionWin = cleanVer(res[1])
	verifyCorrectVersionMust(programVersionWin)
	fmt.Printf("programVersionWin: %s\n", programVersionWin)
}

func extractVersionMacMust() {
	path := filepath.Join("mac", "dbworkbench", "Info.plist")
	lines, err := u.ReadLinesFromFile(path)
	fataliferr(err)
	// extract from:
	// 	<key>CFBundleShortVersionString</key>
	//  <string>0.1</string>
	idx := -1
	for i, l := range lines {
		l = strings.ToLower(strings.TrimSpace(l))
		if strings.Contains(l, "<key>cfbundleshortversionstring</key>") {
			idx = i
			break
		}
	}
	fatalif(idx == -1, "didn't find <key>CFBundleShortVersionString</key>")
	s := strings.TrimSpace(lines[idx+1])
	if strings.HasPrefix(s, "<string>") {
		s = s[len("<string>"):]
	} else {
		fatalf("invalid s: '%s'\n", s)
	}
	if strings.HasSuffix(s, "</string>") {
		s = s[:len(s)-len("</string>")]
	} else {
		fatalf("invalid s: '%s'\n", s)
	}
	verifyCorrectVersionMust(s)
	programVersionMac = s
	fmt.Printf("programVersionMac: %s\n", programVersionMac)
}

func extractVersionMust() {
	extractVersionWinMust()
	extractVersionMacMust()
	fatalif(programVersionMac != programVersionWin, "programVersionMac != programVersionWin ('%s' != '%s')", programVersionMac, programVersionWin)
	programVersion = programVersionMac
	fmt.Printf("programVersion: %s\n", programVersion)
}

func buildMac() {
	dirToZip := filepath.Join("mac", "build", "Release", "Database Workbench.app")
	zipPath := filepath.Join("mac", "build", "Release", "dbworkbenchmacapp.zip")
	err := ZipDirectory(dirToZip, zipPath)
	fataliferr(err)

	// Maybe TODO: Upload to S3
}

func main() {
	parseCmdLine()
	if flgGenResources {
		genResources()
		return
	}

	if isWin() {
		fmt.Printf("Starting windows build\n")
	} else if isMac() {
		fmt.Printf("Starting mac build\n")
	} else {
		log.Fatalf("unsupported os, runtime.GOOS='%s'\n", runtime.GOOS)
	}
	verifyGitCleanMust()
	extractVersionMust()
	if isWin() {
		buildWinAll()
	}

	if isMac() {
		buildMac()
	}
}
