package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kjk/u"
	"github.com/mitchellh/go-homedir"
)

// Options defines cmd-line and computed configuration options
type Options struct {
	// options that come from command-line
	Debug    bool
	URL      string
	Host     string
	Port     int
	User     string
	Pass     string
	DbName   string
	Ssl      string
	HTTPHost string
	HTTPPort int
	IsDev    bool
	Test     bool

	// computed options
	ResourcesFromZip bool
}

var options Options

func parseCmdLine() {
	flag.BoolVar(&options.Debug, "debug", false, "enable debug mode")
	flag.StringVar(&options.URL, "url", "", "database connection string")
	flag.StringVar(&options.Host, "host", "", "database host name or ip address")
	flag.IntVar(&options.Port, "port", 5432, "database port")
	flag.StringVar(&options.User, "user", "", "database user")
	flag.StringVar(&options.Pass, "pass", "", "database password for user")
	flag.StringVar(&options.DbName, "db", "", "database name")
	flag.StringVar(&options.Ssl, "ssl", "", "SSL options")
	// using 127.0.0.1 so that windows firewall doesn't complain about
	// opening externally-accessible ports
	if isWindows() {
		flag.StringVar(&options.HTTPHost, "bind", "127.0.0.1", "HTTP server host")
	} else {
		flag.StringVar(&options.HTTPHost, "bind", "", "HTTP server host")
	}
	flag.IntVar(&options.HTTPPort, "listen", 5444, "HTTP server listen port")
	flag.BoolVar(&options.IsDev, "dev", false, "true for running in dev mode")
	flag.BoolVar(&options.Test, "test", false, "if true, runs a test (whatever it might be)")
	flag.Parse()
}

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func verifyDirs() {
	logDir := getLogDir()
	os.MkdirAll(logDir, 0755)
	if !u.PathExists(logDir) {
		log.Fatalf("directory '%s' doesn't exist. Please create!\n", logDir)
	}
	if !u.PathExists(getDataDir()) {
		log.Fatalf("directory '%s' doesn't exist\n", getDataDir())
	}
}

func getDataDirMac() string {
	d, err := homedir.Dir()
	if err != nil {
		log.Fatalf("homedir.Dir() failed with %s", err)
	}
	d = filepath.Join(d, "Library", "Application Support", "dbHero")
	return d
}

func getDataDirWindows() string {
	dir := os.Getenv("LOCALAPPDATA")
	if dir == "" {
		log.Fatalf("LOCALAPPDATA not set")
	}
	return filepath.Join(dir, "dbhero")
}

func getDataDir() string {
	if isMac() {
		return getDataDirMac()
	}

	if isWindows() {
		return getDataDirWindows()
	}
	log.Fatalf("unsupported runtime.GOOS value: '%s", runtime.GOOS)
	return ""
}

func getLogDir() string {
	return filepath.Join(getDataDir(), "log")
}

func runGulpAndWaitExit() {
	gulpPath := filepath.Join("node_modules", ".bin", "gulp")
	cmd := exec.Command(gulpPath, "build_and_watch")
	cmdStr := strings.Join(cmd.Args, " ")
	fmt.Printf("starting '%s'\n", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start('%s') failed with '%s'\n", cmdStr, err)
	}
	cmd.Wait()
	LogInfof("gulp exited\n")
}

func runGulpAsync() {
	go func() {
		// on error gulp exists which means after first error in JavaScript
		// we stop re-generating it. We need to restart it automatically with
		// few seconds delay to allow for fixing the bug
		for {
			runGulpAndWaitExit()
			time.Sleep(time.Second * 5)
		}
	}()
}

func main() {
	parseCmdLine()

	if options.Test {
		testMysqlInfo()
		os.Exit(0)
	}

	if options.IsDev {
		logToStdout = true
	}
	verifyDirs()
	RemoveOldLogFiles()
	OpenLogFiles()
	IncLogVerbosity()
	LogInfof("Data dir: %s\n", getDataDir())
	if !options.IsDev {
		options.ResourcesFromZip = true
		LogInfof("reading resources from zip\n")
	}

	openUsageFileMust()

	if options.IsDev {
		runGulpAsync()
	}

	if options.ResourcesFromZip {
		err := loadResourcesFromEmbeddedZip()
		if err != nil {
			LogFatalf("loadResourcesFromEmbeddedZip() failed with '%s'\n", err)
		}
	}

	go startWebServer()
	handleSignals()

}
