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

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func isWindows() bool {
	return runtime.GOOS == "windows"
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

func startGulp() {
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
}

func testMysqlInfo() {
	// ${user}:${pwd}@tcp(${ip}:${port})/{dbName}?parseTime=true
	uri := os.Getenv("MYSQL_URL")
	if uri == "" {
		log.Fatal("MYSQL_URL env variable is not defined. in bash, do:\nMYSQL_URL=${uri} ./scripts.run.sh\n")
	}
	c, err := NewClientMysqlFromURL(uri)
	if err != nil {
		fmt.Printf("NewClientMysqlFromURL('%s') failed with '%s'\n", uri, err)
		return
	}
	defer c.Close()
	if err = c.Test(); err != nil {
		fmt.Printf("c.Test() failed with '%s', uri: '%s'\n", err, uri)
		return
	}

	i, err := c.Info()
	if err != nil {
		fmt.Printf("c.Info() failed with '%s'\n", err)
		return
	}
	i.Dump()

	dbs, err := c.Databases()
	if err != nil {
		fmt.Printf("c.Databases() failed with '%s'\n", err)
		return
	}
	fmt.Printf("%d datbases:\n", len(dbs))
	for _, db := range dbs {
		fmt.Printf("  '%s'\n", db)
	}

	tables, err := c.Tables()
	if err != nil {
		fmt.Printf("c.Tables() failed with '%s'\n", err)
		return
	}
	fmt.Printf("%d tables:\n", len(tables))
	for _, s := range tables {
		fmt.Printf("  '%s'\n", s)
	}

	for _, t := range tables {
		schema, err := c.Table(t)
		if err != nil {
			fmt.Printf("c.TableRows('%s') failed with '%s'\n", t, err)
			return
		}
		fmt.Printf("Table: '%s'\n", t)
		schema.DumpFull()
	}
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
		startGulp()
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
