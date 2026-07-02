package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Date       string   `json:"date"`
	Step       int      `json:"step"`
	EntryTime  string   `json:"entry_time,omitempty"`
	LunchStart string   `json:"lunch_start,omitempty"`
	LunchEnd   string   `json:"lunch_end,omitempty"`
	ExitTime   string   `json:"exit_time,omitempty"`
	Logs       []string `json:"logs"`
}

const (
	layout          = "2006-01-02"
	timestampLayout = time.RFC3339
)

func Today() string {
	return time.Now().Format(layout)
}

func defaultState() State {
	return State{Date: Today(), Step: 0}
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
		if home == "" {
			return "", err
		}
	}

	dir := filepath.Join(home, ".worklog")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "state.json"), nil
}

func Load() (State, error) {
	p, err := path()
	if err != nil {
		return defaultState(), err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultState(), nil
		}
		return defaultState(), err
	}

	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return defaultState(), err
	}

	if s.Date != Today() {
		return defaultState(), nil
	}

	return s, nil
}

func Reset() error {
	return Save(defaultState())
}

func parseStamp(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	t, err := time.Parse(timestampLayout, value)
	return t, err == nil
}

func (s State) WorkedDuration() time.Duration {
	entry, ok := parseStamp(s.EntryTime)
	if !ok {
		return 0
	}

	now := time.Now()
	duration := time.Duration(0)

	if start, ok := parseStamp(s.LunchStart); ok {
		duration += start.Sub(entry)

		if end, ok := parseStamp(s.LunchEnd); ok {
			if exit, ok := parseStamp(s.ExitTime); ok {
				duration += exit.Sub(end)
			} else {
				duration += now.Sub(end)
			}
		}
	} else if exit, ok := parseStamp(s.ExitTime); ok {
		duration += exit.Sub(entry)
	} else {
		duration += now.Sub(entry)
	}

	if duration < 0 {
		return 0
	}

	return duration
}

func Save(s State) error {
	s.Date = Today()

	p, err := path()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p, b, 0644)
}
