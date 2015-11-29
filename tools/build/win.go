package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func s3ExeSetupPath() string {
	return s3Dir + fmt.Sprintf("rel/DatabaseWorkbench-setup-%s.exe", programVersion)
}

func exeSetupPath() string {
	exeName := fmt.Sprintf("DatabaseWorkbench-setup-%s.exe", programVersion)
	return pj("bin", "Release", exeName)
}

func exeSetupTmpPath() string {
	return pj("bin", "Release", "DatabaseWorkbench-setup-inno.exe")
}

func exePath() string {
	return pj("bin", "Release", "DatabaseWorkbench.exe")
}

func s3VerifyExists(s3Path string) {
	fatalif(!s3Exists(s3Path), "'%s' doesn't exist in s3\n", s3Path)
}

func s3VerifyNotExists(s3Path string) {
	fatalif(s3Exists(s3Path), "'%s' already exist in s3\n", s3Path)
}

func signMust(path string) {
	// signtool is finicky so we copy cert.pfx to the directory where the file is
	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)
	certPwd := cachedSecrets.CertPwd
	certDest := pj(fileDir, "cert.pfx")
	fileCopyMust(certDest, certPath)
	cmd := getCmdInEnv(getEnvForVS(), "signtool.exe", "sign", "/t", "http://timestamp.verisign.com/scripts/timstamp.dll",
		"/du", "http://databaseworkbench.com", "/f", "cert.pfx",
		"/p", certPwd, fileName)
	cmd.Dir = fileDir
	runCmdMust(cmd, true)
}

func clean() {
	removeDirMust("obj")
	removeDirMust("bin")
}

func winDir() string {
	return filepath.Join("win", "DatabaseWorkbench")
}

func cdToWinDir() {
	err := os.Chdir(winDir())
	fataliferr(err)
}

func copyDbWorkbench() {
	dst := filepath.Join(winDir(), "dbworkbench.exe")
	src := "dbworkbench.exe"
	fileCopyMust(dst, src)

	dst = filepath.Join(winDir(), "dbworkbench.dat")
	src = "dbworkbench.dat"
	fileCopyMust(dst, src)
}

func buildWin() {
	if flgUpload {
		s3VerifyNotExists(s3ExeSetupPath())
	}
	copyDbWorkbench()
	cdToWinDir()
	clean()

	out, err := runMsbuildGetOutput(true, "DatabaseWorkbench.csproj", "/t:Rebuild", "/p:Configuration=Release", "/m")
	if err != nil {
		fmt.Printf("failed with:\n%s\n", string(out))
	}
	fataliferr(err)
}

func buildSetupWin() {
	signMust(exePath())
	signMust("dbworkbench.exe")

	ver := fmt.Sprintf("/dMyAppVersion=%s", programVersion)
	cmd := exec.Command(innoSetupPath, "/Qp", ver, "installer.iss")
	fmt.Printf("Running %s\n", cmd.Args)
	runCmdMust(cmd, true)
	signMust(exeSetupTmpPath())
	fileCopyMust(exeSetupPath(), exeSetupTmpPath())
}

func uploadToS3Win() {
	if !flgUpload {
		return
	}
	s3UploadFile(s3ExeSetupPath(), exeSetupPath(), true)
	s3Url := "https://kjkpub.s3.amazonaws.com/" + s3ExeSetupPath()
	buildOn := time.Now().Format("2006-01-02")
	jsTxt := fmt.Sprintf(`var LatestVerWin = "%s";
var LatestUrlWin = "%s";
var BuiltOnWin = "%s";
`, programVersion, s3Url, buildOn)
	s3UploadString(s3Dir+"latestverwin.js", jsTxt, true)
}

func buildWinAll() {
	verifyHasSecretsMust()
	detectInnoSetupMust()
	buildWin()
	buildSetupWin()
	uploadToS3Win()
}
