package actions

import (
	"fmt"
	"time"

	"worklog/internal/state"
)

func AddLog(s *state.State, msg string) {
	entry := fmt.Sprintf("%s → %s", time.Now().Format("15:04:05"), msg)
	s.Logs = append(s.Logs, entry)
}

func Chkin(s *state.State) {
	switch s.Step {
	case 0:
		AddLog(s, "CHECKIN: Entry")
		s.Step = 1
	case 2:
		AddLog(s, "CHECKIN: Lunch End")
		s.Step = 3
	default:
		fmt.Println("invalid chkin")
	}
}

func Chkout(s *state.State) {
	switch s.Step {
	case 1:
		AddLog(s, "CHKOUT: Lunch Start")
		s.Step = 2
	case 3:
		AddLog(s, "CHKOUT: Exit")
		s.Step = 4
	default:
		fmt.Println("invalid chkout")
	}
}
