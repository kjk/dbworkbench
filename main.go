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
	_ "github.com/lib/pq"
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
	SkipOpen bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	IsLocal  bool   `long:"local" description:"is true if running locally (dev mode)"`
}

var dbClient *Client
var options Options

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func initClient() {
	if connectionSettingsBlank(options) {
		return
	}

	client, err := NewClient()
	if err != nil {
		exitWithMessage(err.Error())
	}

	if options.Debug {
		fmt.Println("Server connection string:", client.connectionString)
	}

	fmt.Println("Connecting to server...")
	err = client.Test()
	if err != nil {
		exitWithMessage(err.Error())
	}

	fmt.Println("Checking tables...")
	_, err = client.Tables()
	if err != nil {
		exitWithMessage(err.Error())
	}

	dbClient = client
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

func openPage() {
	if options.SkipOpen {
		return
	}

	url := fmt.Sprintf("http://%v:%v", options.HTTPHost, options.HTTPPort)
	fmt.Println("To view database open", url, "in browser")

	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	exec.Command("open", url).Output()
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
	if options.IsLocal {
		return u.ExpandTildeInPath("~/data/dbworkbench")
	}
	//  on the server it's in /home/dbworkbench/www/data
	return u.ExpandTildeInPath("~/www/data")
}

func getLogDir() string {
	return filepath.Join(getDataDir(), "log")
}

func startWebpackWatch() {
	cmd := exec.Command("./scripts/webpack-dev.sh")
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
	runtime.GOMAXPROCS(runtime.NumCPU())

	initOptions()

	logToStdout = true
	verifyDirs()
	OpenLogFiles()
	IncLogVerbosity()
	LogInfof("local: %v, data dir: %s\n", options.IsLocal, getDataDir())
	initCookieMust()

	getDbMust()

	if options.IsLocal {
		startWebpackWatch()
	}

	initClient()

	if dbClient != nil {
		defer dbClient.db.Close()
	}

	if options.Debug {
		startRuntimeProfiler()
	}

	go startWebServer()
	openPage()
	handleSignals()
}
