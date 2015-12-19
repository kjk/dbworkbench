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
	"time"

	"github.com/kjk/u"
)

const (
	s3Dir = "software/dbhero/"
)

var (
	flgNoCleanCheck     bool
	flgUpload           bool
	flgUploadAutoUpdate bool
	flgGenResources     bool
	flgBeta             bool

	programVersionWin      string
	programVersionMac      string
	programVersionFrontend string
	programVersion         string
	certPath               string
	cachedSecrets          *Secrets
	innoSetupPath          string
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
	flag.BoolVar(&flgBeta, "beta", false, "if ture, build a beta version")
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

func verifyHasCertWinMust() {
	if !isWin() {
		return
	}
	certPath = pj("scripts", "cert.pfx")
	if !fileExists(certPath) {
		certPath = pj("..", "..", "..", "..", "..", "sumatrapdf", "scripts", "cert.pfx")
		fatalif(!fileExists(certPath), "didn't find cert.pfx in scripts/ or ../../../sumatrapdf/scripts/\n")
	}
	absPath, err := filepath.Abs(certPath)
	fataliferr(err)
	certPath = absPath
	fmt.Printf("signing certificate: '%s'\n", certPath)
}

func verifyHasSecretsMust() {
	verifyHasCertWinMust()
	secretsPath := pj("scripts", "secrets.json")
	if !fileExists(secretsPath) {
		secretsPath = pj("..", "..", "..", "..", "..", "sumatrapdf", "scripts", "secrets.json")
		fatalif(!fileExists(secretsPath), "didn't find secrets.json in scripts/ or ../../../sumatrapdf/scripts/\n")
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
	path := filepath.Join("win", "dbhero", "Properties", "AssemblyInfo.cs")
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

/* given a text in .plist:
<key>foo<key>
<string>bar</key>
returns "bar" if keyName is "foo"
*/
func plistGetStrVal(lines []string, keyName string) string {
	idx := -1
	key := fmt.Sprintf("<key>%s</key>", strings.ToLower(keyName))
	for i, l := range lines {
		l = strings.ToLower(strings.TrimSpace(l))
		if strings.Contains(l, key) {
			idx = i
			break
		}
	}
	fatalif(idx == -1, "didn't find <key>%s</key>\n", keyName)
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
	return cleanVer(s)
}

func extractVersionMacMust() {
	path := filepath.Join("mac", "dbHero", "Info.plist")
	lines, err := u.ReadLinesFromFile(path)
	fataliferr(err)
	// extract from:
	// 	<key>CFBundleShortVersionString</key>
	//  <string>0.1</string>
	shortVer := plistGetStrVal(lines, "CFBundleShortVersionString")
	ver := plistGetStrVal(lines, "CFBundleVersion")
	verifyCorrectVersionMust(shortVer)
	verifyCorrectVersionMust(ver)
	fatalif(shortVer != ver, "shortVer (%s) != ver (%s)\n", shortVer, ver)
	programVersionMac = shortVer
	fmt.Printf("programVersionMac: %s\n", programVersionMac)
}

func extractVersionFrontendMust() {
	path := filepath.Join("s", "index.html")
	d, err := ioutil.ReadFile(path)
	fataliferr(err)
	r := regexp.MustCompile(`var gVersionNumber = "(.*)"`)
	s := string(d)
	res := r.FindStringSubmatch(s)

	programVersionFrontend = cleanVer(res[1])
	verifyCorrectVersionMust(programVersionFrontend)
	fmt.Printf("programVersionFrontend: %s\n", programVersionFrontend)
}

func extractVersionMust() {
	extractVersionWinMust()
	extractVersionMacMust()
	extractVersionFrontendMust()
	fatalif(programVersionMac != programVersionWin, "programVersionMac != programVersionWin ('%s' != '%s')\n", programVersionMac, programVersionWin)
	fatalif(programVersionMac != programVersionFrontend, "programVersionMac != programVersionFrontend ('%s' != '%s')\n", programVersionMac, programVersionFrontend)
	programVersion = programVersionMac
	fmt.Printf("programVersion: %s\n", programVersion)
}

func s3SetupPathMac() string {
	return s3Dir + fmt.Sprintf("rel/dbHero-%s.zip", programVersion)
}

func s3SetupPathMacBeta() string {
	return s3Dir + fmt.Sprintf("beta/dbHero-%s.zip", programVersion)
}

func macZipPath() string {
	return pj("mac", "build", "Release", "dbHero.zip")
}

func uploadToS3Mac() {
	if !flgUpload {
		fmt.Printf("skipping s3 upload because -upload not given\n")
		return
	}

	verifyHasSecretsMust()

	s3Path := s3SetupPathMac()
	s3VerifyNotExistsMust(s3Path)

	s3UploadFile(s3Path, macZipPath(), true)
	s3Url := "https://kjkpub.s3.amazonaws.com/" + s3Path
	buildOn := time.Now().Format("2006-01-02")
	jsTxt := fmt.Sprintf(`var LatestVerMac = "%s";
var LatestUrlMac = "%s";
var BuiltOnMac = "%s";
`, programVersion, s3Url, buildOn)
	s3UploadString(s3Dir+"latestvermac.js", jsTxt, true)
	s3VerifyExistsWaitMust(s3Path)
}

func uploadToS3MacBeta() {
	if !flgUpload {
		fmt.Printf("skipping s3 upload because -upload not given\n")
		return
	}

	verifyHasSecretsMust()

	s3Path := s3SetupPathMacBeta()
	s3VerifyNotExistsMust(s3Path)

	s3UploadFile(s3Path, macZipPath(), true)
	s3Url := "https://kjkpub.s3.amazonaws.com/" + s3Path
	buildOn := time.Now().Format("2006-01-02")
	jsTxt := fmt.Sprintf(`var LatestVerMac = "%s";
var LatestUrlMac = "%s";
var BuiltOnMac = "%s";
`, programVersion, s3Url, buildOn)
	s3UploadString(s3Dir+"latestvermacbeta.js", jsTxt, true)
	s3VerifyExistsWaitMust(s3Path)
}

func buildMac() {
	dirToZip := filepath.Join("mac", "build", "Release", "dbHero.app")
	zipPath := filepath.Join("mac", "build", "Release", "dbHero.zip")
	err := ZipDirectory(dirToZip, zipPath)
	fataliferr(err)

	if flgBeta {
		uploadToS3MacBeta()

	} else {
		uploadToS3Mac()
	}
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
	} else if isMac() {
		buildMac()
	} else {
		log.Fatalf("unsupported runtime.GOOS: '%s'\n", runtime.GOOS)
	}
}
