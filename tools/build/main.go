package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	flgNoCleanCheck bool
	flgUpload       bool

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

func detectInnoSetupMust() {
	path1 := pj(os.Getenv("ProgramFiles"), "Inno Setup 5", "iscc.exe")
	if fileExists(path1) {
		innoSetupPath = path1
		fmt.Printf("Inno Setup: %s\n", innoSetupPath)
		return
	}
	path2 := pj(os.Getenv("ProgramFiles(x86)"), "Inno Setup 5", "iscc.exe")
	if fileExists(path2) {
		innoSetupPath = path2
		fmt.Printf("Inno Setup: %s\n", innoSetupPath)
		return
	}
	fatalif(true, "didn't find Inno Setup (tried '%s' and '%s'). Download from http://www.jrsoftware.org/isinfo.php\n", path1, path2)
}

func parseCmdLine() {
	// -no-clean-check is useful when testing changes to this build script
	flag.BoolVar(&flgNoCleanCheck, "no-clean-check", false, "allow running if repo has changes (for testing build script)")
	flag.BoolVar(&flgUpload, "upload", false, "if true, upload the files")
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
	// TODO: write me
	programVersionMac = programVersionWin

	verifyCorrectVersionMust(programVersionMac)
	fmt.Printf("programVersionMac: %s\n", programVersionMac)
}

func extractVersionMust() {
	extractVersionWinMust()
	extractVersionMacMust()
	fatalif(programVersionMac != programVersionWin, "programVersionMac != programVersionWin ('%s' != '%s')", programVersionMac, programVersionWin)
	programVersion = programVersionMac
	fmt.Printf("programVersion: %s\n", programVersion)
}

func main() {
	parseCmdLine()
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
		detectInnoSetupMust()
	}

}
