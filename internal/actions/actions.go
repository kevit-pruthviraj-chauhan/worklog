package actions

import (
	"errors"
	"fmt"
	"time"

	"worklog/internal/state"
)

func AddLog(s *state.State, msg string, t time.Time) {
	entry := fmt.Sprintf("%s → %s", t.Format("15:04:05"), msg)
	s.Logs = append(s.Logs, entry)
}

func nowLocal() time.Time {
	now := time.Now()
	return now.Local()
}

func ParseOptionalTime(arg string, ref time.Time) (time.Time, error) {
	if arg == "" {
		return nowLocal(), nil
	}

	parsed, err := time.Parse("15:04", arg)
	if err != nil {
		return time.Time{}, err
	}
	now := nowLocal()
	result := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
	if !ref.IsZero() && result.Before(ref) {
		candidate := result.Add(12 * time.Hour)
		if !candidate.Before(ref) && candidate.Sub(ref) < 24*time.Hour {
			return candidate, nil
		}
	}

	return result, nil
}

func Chkin(s *state.State, ts time.Time) error {
	switch s.Step {
	case 0:
		AddLog(s, "CHECKIN: Entry", ts)
		s.EntryTime = ts.Format(time.RFC3339)
		s.Step = 1
		return nil
	case 2:
		AddLog(s, "CHECKIN: Lunch End", ts)
		s.LunchEnd = ts.Format(time.RFC3339)
		s.Step = 3
		return nil
	default:
		return errors.New("cannot check in right now")
	}
}

func Chkout(s *state.State, ts time.Time) error {
	switch s.Step {
	case 1:
		AddLog(s, "CHKOUT: Lunch Start", ts)
		s.LunchStart = ts.Format(time.RFC3339)
		s.Step = 2
		return nil
	case 3:
		AddLog(s, "CHKOUT: Exit", ts)
		s.ExitTime = ts.Format(time.RFC3339)
		s.Step = 4
		return nil
	default:
		return errors.New("cannot check out right now")
	}
}
