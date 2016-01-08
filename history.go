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

// GetHistory returns history records, most recent at the beginning
func (h *History) GetHistory() []HistoryRecord {
	n := len(h.history)
	res := make([]HistoryRecord, 0)
	for i := n - 1; i >= 0; i-- {
		res = append(res, h.history[i])
	}
	return res
}

// AddToHistory remembers query in history
func (h *History) AddToHistory(query string) {
	hr := HistoryRecord{
		Query:     query,
		Timestamp: time.Now().String(),
	}
	h.history = append(h.history, hr)
}
