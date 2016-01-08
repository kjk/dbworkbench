package main

import (
	"time"
)

// HistoryRecord remembers a single history itema
type HistoryRecord struct {
	Query     string `json:"query"`
	Timestamp string `json:"timestamp"`
}

// History remembers history records
type History struct {
	history []HistoryRecord
}

// History returns history records
func (h *History) GetHistory() []HistoryRecord {
	return h.history
}

// AddToHistory remembers query in history
func (h *History) AddToHistory(query string) {
	hr := HistoryRecord{
		Query:     query,
		Timestamp: time.Now().String(),
	}
	h.history = append(h.history, hr)
}
