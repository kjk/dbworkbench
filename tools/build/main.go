package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	flgNoCleanCheck bool
	flgUpload       bool
	flgWin          bool
	flgMac          bool
)

func finalizeThings(crashed bool) {
	// nothing for us
}

func parseCmdLine() {
	// -no-clean-check is useful when testing changes to this build script
	flag.BoolVar(&flgNoCleanCheck, "no-clean-check", false, "allow running if repo has changes (for testing build script)")
	flag.BoolVar(&flgWin, "win", false, "windows build")
	flag.BoolVar(&flgMac, "mac", false, "mac build")
	flag.Parse()
	fatalif(!flgWin && !flgMac, "-win or -mac cmd-line argument must be given")
	fatalif(flgWin && flgMac, "don't specify both -win and -mac")
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

func main() {
	fmt.Printf("Starting windows build\n")
	parseCmdLine()
}
