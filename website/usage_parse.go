package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

// parsing of usage.txt, stats based on that

const (
	// OsMac means running on Mac
	OsMac = 1
	// OsWindows running on Windows
	OsWindows = 2
)

func testUsageParse() {
	path := "usage-for-testing.go"
	parseUsage(path)
}

// UsageUser describes statistics about users
type UsageUser struct {
	ID                string
	Name              string
	UniqueDaysCount   int
	FirstSeen         time.Time
	LastSeen          time.Time
	DatabaseOpenCount int
	QueryCount        int
	Os                int // 1 : mac, 2 : windows
	OsVersion         string
}

// records are separated by single empty line
func readUsageRecord(r *bufio.Reader) ([]string, error) {
	var res []string
	var err error
	var l string

	for {
		l, err = r.ReadString('\n')
		// result of ReadString() includes terminating \n, so remove it
		l = strings.TrimRight(l, "\n")
		if l == "" {
			break
		}
	}
	return res, err
}

// UsageRecord describes data from a single usage
type UsageRecord struct {
	ID                string
	Name              string
	When              time.Time
	Os                int
	Version           string
	OsVersion         string
	DatabaseOpenCount int
	QueryCount        int
}

func parseUsageRecord(lines []string) (*UsageRecord, error) {
	var res UsageRecord
	return &res, nil
}

func parseUsage(path string) ([]*UsageUser, error) {
	var res []*UsageUser
	file, err := os.Open(path)
	if err != nil {
		LogErrorf("os.Open() failed with '%'s\n", err)
		return nil, err
	}
	defer file.Close()

	idToUser := make(map[string]*UsageUser)

	r := bufio.NewReader(file)
	for {
		lines, err := readUsageRecord(r)
		if err != nil && err != io.EOF {
			return res, err
		}
		if len(lines) == 0 {
			return res, nil
		}
		rec, err := parseUsageRecord(lines)
		if err != nil {
			LogErrorf("parseUsageRecord() of '%v' failed with '%s'\n", lines, err)
			continue
		}
		user := idToUser[rec.ID]
		if user == nil {
			user = &UsageUser{}
			user.ID = rec.ID
			user.Name = rec.Name
			user.Os = rec.Os
			idToUser[rec.ID] = user
			user.FirstSeen = rec.When
		}
		// remember the latest version
		user.OsVersion = rec.OsVersion
		user.LastSeen = rec.When
	}
}
