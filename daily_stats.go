package main

import (
	"encoding/json"
	"sync"
	"time"
)

/*
We want a summary of important stats to be gathered daily and e-mailed to us.

Settings are also saved in ${data}/stats/YYYY-MM/DD.json file so that
we can analyze them in the future.

TODO:
 - more stats
 - graceful restarts (save on signal.SIGQUIT, reaload on startup)
*/

var (
	muDailyStats  sync.Mutex
	dailyStatsDay time.Time
	dailyStats    *DailyStats
)

// DailyStats contains info about the most important stats for a given day
type DailyStats struct {
	// user ids of new user accounts
	NewUserAccounts []int
	// user ids of users that logged in
	LoggedUsers []int
}

func saveDailyStats(d []byte, day time.Time) {
	// TODO: write me
}

func reloadCurrentDayOnStartup() {
	// TODO: write me
}

func rotateDailyStatsIfNecessary() {
	if dailyStats == nil {
		dailyStatsDay = time.Now()
		dailyStats = &DailyStats{}
		return
	}
	day := time.Now().YearDay()
	currDay := dailyStatsDay.YearDay()
	if day == currDay {
		return
	}
	d, err := json.Marshal(dailyStats)
	if err != nil {
		LogErrorf("json.Marshal() failed with '%s'\n", err)
		go saveDailyStats(d, dailyStatsDay)
	}
	dailyStatsDay = time.Now()
	dailyStats = &DailyStats{}
}

func withDailyStatsLocked(f func(*DailyStats)) {
	muDailyStats.Lock()
	rotateDailyStatsIfNecessary()
	f(dailyStats)
	muDailyStats.Unlock()
}

func dsAccountCreate(userID int) {
	withDailyStatsLocked(func(ds *DailyStats) {
		ds.NewUserAccounts = IntAppendIfNotExists(ds.NewUserAccounts, userID)
	})
}

func dsUserLogin(userID int) {
	withDailyStatsLocked(func(ds *DailyStats) {
		ds.LoggedUsers = IntAppendIfNotExists(ds.LoggedUsers, userID)
	})
}
