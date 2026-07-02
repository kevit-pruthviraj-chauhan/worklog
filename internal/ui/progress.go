package ui

import (
	"fmt"
	"strings"

	"worklog/internal/state"
)

func Progress(s state.State) {
	total := 4
	done := s.Step

	barLen := 30
	filled := (done * barLen) / total

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)

	fmt.Println("\nTODAY PROGRESS")
	fmt.Printf("[%s] %d/%d\n\n", bar, done, total)

	fmt.Println("STATE:", label(s.Step))
}

func label(step int) string {
	switch step {
	case 0:
		return "Not started"
	case 1:
		return "At work (morning)"
	case 2:
		return "Lunch break"
	case 3:
		return "Back from lunch"
	case 4:
		return "Day complete"
	default:
		return "unknown"
	}
}