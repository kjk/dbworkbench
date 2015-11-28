package main

import (
	"log"
	"os"
	"runtime"
	"time"
)

const megabyte = 1024 * 1024

func startRuntimeProfiler() {
	go func() {
		logger := log.New(os.Stdout, "", 0)
		m := &runtime.MemStats{}

		for {
			runtime.ReadMemStats(m)

			logger.Printf(
				"[DEBUG] Goroutines: %v, Mem used: %v (%v mb), Mem acquired: %v (%v mb)\n",
				runtime.NumGoroutine(),
				m.Alloc, m.Alloc/megabyte,
				m.Sys, m.Sys/megabyte,
			)

			time.Sleep(time.Second * 30)
		}
	}()
}
