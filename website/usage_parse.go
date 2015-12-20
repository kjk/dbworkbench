package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// parsing of usage.txt, stats based on that

const (
	// OsUnknown is default value for when os type is missing
	OsUnknown = 0
	// OsMac means running on Mac
	OsMac = 1
	// OsWindows running on Windows
	OsWindows = 2
)

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
	Version           string

	// values only used during stats calculations
	days []time.Time // list of unique days
}

func (u *UsageUser) addDayIfUnique(t time.Time) {
	for _, t2 := range u.days {
		if sameDay(t, t2) {
			return
		}
	}
	u.days = append(u.days, t)
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
		//fmt.Printf("l: '%s'\n", l)
		res = append(res, l)
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

// return true if t1 and t2 represent the same day
func sameDay(t1, t2 time.Time) bool {
	// TODO: write me
	return false
}

func parseUsageRecord(lines []string) (*UsageRecord, error) {
	var res UsageRecord
	// "key: value" lines
	for len(lines) > 0 {
		s := lines[0]
		lines = lines[1:]
		if strings.HasPrefix(s, "-----") {
			break
		}
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			LogErrorf("invalid line: '%s'\n", s)
		}
		val := strings.TrimSpace(parts[1])
		switch parts[0] {
		case "user":
			res.Name = val
		case "time":
			secs, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				LogErrorf("invalid time value '%s' in line '%s', err: '%s'\n", val, s, err)
			} else {
				res.When = time.Unix(secs, 0)
			}
		case "ver_in_url":
			// "ver" takes precedence
			if res.Version == "" {
				res.Version = val
			}
		case "ver":
			// takes precedence over "ver_in_url"
			res.Version = val
		case "ostype":
			if val == "mac" {
				res.Os = OsMac
			} else if val == "windows" {
				res.Os = OsWindows
			} else {
				LogErrorf("unknown ostype '%s' in line '%s'\n", val, s)
			}
		case "osversion", "os":
			// a bug of sorts: windows sends "os", mac sends "osversion"
			res.OsVersion = val
		case "serial":
			// unique id on mac
			res.ID = val
		case "networkCardId":
			// unique id on windows
			res.ID = val
		case "day", "ip", "machine", "net", "os_in_url":
			// ignore
		default:
			LogErrorf("unknown key '%s' in line '%s'\n", parts[0], s)
		}
	}

	// json lines
	for len(lines) > 0 {
		s := lines[0]
		lines = lines[1:]
		var v interface{}
		err := json.Unmarshal([]byte(s), &v)
		if err != nil {
			LogErrorf("json.Unmarshal() of '%s' failed with '%s'\n")
			continue
		}
		if m, ok := v.(map[string]interface{}); ok {
			event := m["e"]
			switch event {
			case "do":
				res.DatabaseOpenCount++
			case "qe":
				res.QueryCount++
			default:
				LogErrorf("Unknown event '%s' in '%s'\n", event, s)
			}
		} else {
			LogErrorf("json.Unmarshal() produced value of type %T and not map[string]interface{}\n", v)
		}
	}
	return &res, nil
}

func finalizeUsageRecords(arr []*UsageUser) []*UsageUser {
	for _, v := range arr {
		v.UniqueDaysCount = len(v.days)
		v.days = nil
	}
	return arr
}

func parseUsage(path string) ([]*UsageUser, error) {
	var res []*UsageUser
	file, err := os.Open(path)
	if err != nil {
		LogErrorf("os.Open() failed with '%s'\n", err)
		return nil, err
	}
	defer file.Close()

	idToUser := make(map[string]*UsageUser)

	r := bufio.NewReader(file)
	for {
		lines, err := readUsageRecord(r)
		if err != nil && err != io.EOF {
			return finalizeUsageRecords(res), err
		}
		if len(lines) == 0 {
			return finalizeUsageRecords(res), nil
		}
		rec, err := parseUsageRecord(lines)
		if err != nil {
			LogErrorf("parseUsageRecord() of '%v' failed with '%s'\n", lines, err)
			continue
		}
		if rec.ID == "" {
			// early records didn't have an id
			continue
		}
		user := idToUser[rec.ID]
		if user == nil {
			user = &UsageUser{}
			user.ID = rec.ID
			user.Name = rec.Name
			user.Os = rec.Os
			user.FirstSeen = rec.When

			idToUser[rec.ID] = user
			res = append(res, user)
		}
		user.LastSeen = rec.When
		user.DatabaseOpenCount += rec.DatabaseOpenCount
		user.QueryCount += rec.QueryCount
		user.addDayIfUnique(rec.When)
		// remember the latest versions
		user.OsVersion = rec.OsVersion
		user.Version = rec.Version
	}
}

func testUsageParse() {
	path := "usage-for-testing.txt"
	timeStart := time.Now()
	recs, err := parseUsage(path)
	if err != nil {
		LogErrorf("parseUsage() failed with '%s'\n")
		return
	}
	LogInfof("%d users, time to parse: %s\n", len(recs), time.Since(timeStart))
	s, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		LogErrorf("json.MarshalIndent() failed with '%s'\n", err)
	} else {
		fmt.Printf("%s\n", string(s))
	}
}
