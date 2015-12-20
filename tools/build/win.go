package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func s3SetupPathWin() string {
	return s3Dir + fmt.Sprintf("rel/dbHero-setup-%s.exe", programVersion)
}

func s3SetupPathWinBeta() string {
	return s3Dir + fmt.Sprintf("beta/dbHero-setup-%s.exe", programVersion)
}

func exeSetupPath() string {
	exeName := fmt.Sprintf("dbHero-setup-%s.exe", programVersion)
	return pj("bin", "Release", exeName)
}

func exeSetupTmpPath() string {
	return pj("bin", "Release", "dbHero-setup-inno.exe")
}

func exePath() string {
	return pj("bin", "Release", "dbHero.exe")
}

func signMust(path string) {
	// signtool is finicky so we copy cert.pfx to the directory where the file is
	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)
	certPwd := cachedSecrets.CertPwd
	certDest := pj(fileDir, "cert.pfx")
	fileCopyMust(certDest, certPath)
	cmd := getCmdInEnv(getEnvForVS(), "signtool.exe", "sign", "/t", "http://timestamp.verisign.com/scripts/timstamp.dll",
		"/du", "http://dbheroapp.com", "/f", "cert.pfx",
		"/p", certPwd, fileName)
	cmd.Dir = fileDir
	runCmdMust(cmd, true)
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

func cleanWin() {
	removeDirMust("obj")
	removeDirMust("bin")
}

func winDir() string {
	return filepath.Join("win", "dbhero")
}

func cdToWinDir() {
	err := os.Chdir(winDir())
	fataliferr(err)
}

func buildWin() {
	if flgUpload {
		s3VerifyNotExistsMust(s3SetupPathWin())
	}
	cdToWinDir()
	cleanWin()

	out, err := runMsbuildGetOutput(true, "DBHero.csproj", "/t:Rebuild", "/p:Configuration=Release", "/m")
	if err != nil {
		fmt.Printf("failed with:\n%s\n", string(out))
	}
	fataliferr(err)
}

func buildSetupWin() {
	signMust(exePath())
	signMust("dbherohelper.exe")

	ver := fmt.Sprintf("/dMyAppVersion=%s", programVersion)
	cmd := exec.Command(innoSetupPath, "/Qp", ver, "installer.iss")
	fmt.Printf("Running %s\n", cmd.Args)
	runCmdMust(cmd, true)
	signMust(exeSetupTmpPath())
	fileCopyMust(exeSetupPath(), exeSetupTmpPath())
}

func uploadToS3Win() {
	if !flgUpload {
		fmt.Printf("skipping s3 upload because -upload not given\n")
		return
	}

	s3Path := s3SetupPathWin()
	s3VerifyNotExistsMust(s3Path)

	s3UploadFile(s3Path, exeSetupPath(), true)
	s3Url := "https://kjkpub.s3.amazonaws.com/" + s3Path
	buildOn := time.Now().Format("2006-01-02")
	jsTxt := fmt.Sprintf(`var LatestVerWin = "%s";
var LatestUrlWin = "%s";
var BuiltOnWin = "%s";
`, programVersion, s3Url, buildOn)
	s3UploadString(s3Dir+"latestverwin.js", jsTxt, true)
	s3VerifyExistsWaitMust(s3Path)
}

func uploadToS3WinBeta() {
	if !flgUpload {
		fmt.Printf("skipping s3 upload because -upload not given\n")
		return
	}

	s3Path := s3SetupPathWinBeta()
	s3VerifyNotExistsMust(s3Path)

	s3UploadFile(s3Path, exeSetupPath(), true)
	s3Url := "https://kjkpub.s3.amazonaws.com/" + s3Path
	buildOn := time.Now().Format("2006-01-02")
	jsTxt := fmt.Sprintf(`var LatestVerWin = "%s";
var LatestUrlWin = "%s";
var BuiltOnWin = "%s";
`, programVersion, s3Url, buildOn)
	s3UploadString(s3Dir+"latestverwinbeta.js", jsTxt, true)
	s3VerifyExistsWaitMust(s3Path)
}

func buildWinAll() {
	verifyHasSecretsMust()
	detectInnoSetupMust()
	buildWin()
	buildSetupWin()
	if flgBeta {
		uploadToS3WinBeta()
	} else {
		uploadToS3Win()
	}
}
