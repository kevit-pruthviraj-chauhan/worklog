package ui

import (
	"fmt"
	"strings"
	"time"

	"worklog/internal/state"
)

const (
	workdayTarget   = 8*time.Hour + 30*time.Minute
	dashboardWidth  = 56
	progressBarSize = 36
)

func Progress(s state.State) {
	worked := s.WorkedDuration()
	remaining := workdayTarget - worked
	checkoutTime := time.Now().Add(remaining)
	if remaining < 0 {
		remaining = 0
	}

	printLine()
	printCentered("WORKLOG DASHBOARD")
	printLine()
	printField("Date", s.Date)
	printField("Status", label(s.Step))
	printField("Next", nextAction(s.Step))
	printField("Worked", fmt.Sprintf("%s / %s", formatDuration(worked), formatDuration(workdayTarget)))
	if worked < workdayTarget {
		printField("Remaining", formatDuration(remaining))
		printField("Checkout at", checkoutTime.Format("15:04"))
	} else {
		printField("Overtime", formatDuration(worked-workdayTarget))
	}
	printField("Progress", fmt.Sprintf("%s %d%%", progressBar(worked, workdayTarget), int(percent(worked, workdayTarget))))
	printLine()

	fmt.Println("TIMELINE")
	printTimelineEntry("Entry", s.EntryTime)
	printTimelineEntry("Lunch start", s.LunchStart)
	printTimelineEntry("Lunch end", s.LunchEnd)
	printTimelineEntry("Exit", s.ExitTime)
	printLine()
}

func printLine() {
	fmt.Println(strings.Repeat("=", dashboardWidth))
}

func printCentered(value string) {
	innerWidth := dashboardWidth - 2
	padding := innerWidth - len(value)
	if padding < 0 {
		padding = 0
	}
	left := padding / 2
	right := padding - left
	fmt.Printf("=%s%s%s=\n", strings.Repeat(" ", left), value, strings.Repeat(" ", right))
}

func printField(label, value string) {
	fmt.Printf("%-14s %s\n", label+":", value)
}

func printTimelineEntry(label, value string) {
	fmt.Printf("  %-12s %s\n", label+":", formatTime(value))
}

func progressBar(worked, target time.Duration) string {
	ratio := percent(worked, target) / 100.0
	filled := int(ratio * float64(progressBarSize))
	if filled < 0 {
		filled = 0
	}
	if filled > progressBarSize {
		filled = progressBarSize
	}
	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat("░", progressBarSize-filled))
}

func percent(worked, target time.Duration) float64 {
	if target <= 0 {
		return 0
	}
	pct := float64(worked) / float64(target) * 100
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}

func formatTime(value string) string {
	if value == "" {
		return "-"
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "-"
	}
	return t.Format("15:04")
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %02dm", hours, minutes)
}

func label(step int) string {
	switch step {
	case 0:
		return "Not started"
	case 1:
		return "At work"
	case 2:
		return "Lunch break"
	case 3:
		return "Back from lunch"
	case 4:
		return "Day complete"
	default:
		return "Unknown"
	}
}

func nextAction(step int) string {
	switch step {
	case 0:
		return "Check in"
	case 1:
		return "Start lunch"
	case 2:
		return "End lunch"
	case 3:
		return "Check out"
	case 4:
		return "Day complete"
	default:
		return "Unknown"
	}
}
