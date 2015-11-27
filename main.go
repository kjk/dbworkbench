package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/kjk/u"
)

// Options represents command-line options
type Options struct {
	Debug    bool   `short:"d" long:"debug" description:"Enable debugging mode" default:"false"`
	URL      string `long:"url" description:"Database connection string"`
	Host     string `long:"host" description:"Server hostname or IP"`
	Port     int    `long:"port" description:"Server port" default:"5432"`
	User     string `long:"user" description:"Database user"`
	Pass     string `long:"pass" description:"Password for user"`
	DbName   string `long:"db" description:"Database name"`
	Ssl      string `long:"ssl" description:"SSL option"`
	HTTPHost string `long:"bind" description:"HTTP server host" default:"localhost"`
	HTTPPort uint   `long:"listen" description:"HTTP server listen port" default:"5444"`
	AuthUser string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass string `long:"auth-pass" description:"HTTP basic auth password"`
	IsDev    bool   `long:"dev" description:"is true if running in dev mode"`
}

var options Options

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func initOptions() {
	_, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		log.Fatalf("flags.ParseArgs() failed with %s", err)
	}
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func verifyDirs() {
	if !u.PathExists(getLogDir()) {
		log.Fatalf("directory '%s' doesn't exist. Please create! \n", getLogDir())
	}
	if !u.PathExists(getDataDir()) {
		log.Fatalf("directory '%s' doesn't exist\n", getDataDir())
	}
}

func getDataDir() string {
	return u.ExpandTildeInPath("~/data/dbworkbench")
}

func getLogDir() string {
	return filepath.Join(getDataDir(), "log")
}

func startGulpUnix() {
	cmd := exec.Command("./scripts/run_gulp_watch.sh")
	cmdStr := strings.Join(cmd.Args, " ")
	fmt.Printf("starting '%s'\n", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start('%s') failed with '%s'\n", cmdStr, err)
	}
}

func main() {
	fmt.Printf("starting\n")

	initOptions()

	logToStdout = true
	verifyDirs()
	OpenLogFiles()
	IncLogVerbosity()
	LogInfof("Data dir: %s\n", getDataDir())

	if options.IsDev && runtime.GOOS != "windows" {
		startGulpUnix()
	}

	if options.Debug {
		startRuntimeProfiler()
	}

	go startWebServer()
	handleSignals()
}
