package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Date string   `json:"date"`
	Step int      `json:"step"`
	Logs []string `json:"logs"`
}

const layout = "2006-01-02"

func Today() string {
	return time.Now().Format(layout)
}

func path() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".worklog")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "state.json")
}

func Load() State {
	b, err := os.ReadFile(path())
	if err != nil {
		return State{Date: Today(), Step: 0}
	}

	var s State
	json.Unmarshal(b, &s)

	// auto reset daily
	if s.Date != Today() {
		return State{Date: Today(), Step: 0}
	}

	return s
}

func Save(s State) {
	s.Date = Today()
	b, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(path(), b, 0644)
}
