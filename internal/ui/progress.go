package ui

import (
	"fmt"
	"strings"
	"time"

	"worklog/internal/state"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	workdayTarget = 8*time.Hour + 30*time.Minute

	// ANSI Color codes
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Cyan    = "\033[36m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
	Blue    = "\033[34m"
	Red     = "\033[31m"
)

func Progress(s state.State) {
	worked := s.WorkedDuration()
	remaining := workdayTarget - worked
	checkoutTime := time.Now().Add(remaining)
	if remaining < 0 {
		remaining = 0
	}

	fmt.Println()
	printBanner()
	fmt.Println()

	// Date and Status
	fmt.Printf("%s▸ DATE:%s      %s\n", Cyan, Reset, s.Date)
	statusSymbol := getStatusSymbol(s.Step)
	statusColor := getStatusColor(s.Step)
	fmt.Printf("%s▸ STATUS:%s    %s%s%s %s%s\n", Cyan, Reset, statusColor, statusSymbol, Reset, Bold, label(s.Step))
	fmt.Println()

	// Check in and out times
	checkInTime := formatTime(s.EntryTime)
	checkOutTime := formatTime(s.ExitTime)
	fmt.Printf("%s┌─ TIMES ─────────────────────────────┐%s\n", Green, Reset)
	fmt.Printf("%s│%s CHECK IN:  %s%s%s\n", Green, Reset, Green, checkInTime, Reset)
	if checkOutTime == "-" {
		fmt.Printf("%s│%s CHECK OUT: %s%s (est)%s\n", Green, Reset, Yellow, checkoutTime.Format("3:04 PM"), Reset)
	} else {
		fmt.Printf("%s│%s CHECK OUT: %s%s%s\n", Green, Reset, Green, checkOutTime, Reset)
	}
	fmt.Printf("%s└─────────────────────────────────────┘%s\n", Green, Reset)
	fmt.Println()

	// Lunch times
	lunchStart := formatTime(s.LunchStart)
	lunchEnd := formatTime(s.LunchEnd)
	fmt.Printf("%s┌─ LUNCH ─────────────────────────────┐%s\n", Blue, Reset)
	fmt.Printf("%s│%s START: %s%s   END: %s%s%s\n", Blue, Reset, Blue, lunchStart, Blue, lunchEnd, Reset)
	fmt.Printf("%s└─────────────────────────────────────┘%s\n", Blue, Reset)
	fmt.Println()

	// Work hours with visual indicator
	percentWorked := int(percent(worked, workdayTarget))
	fmt.Printf("%s┌─ HOURS ─────────────────────────────┐%s\n", Yellow, Reset)
	fmt.Printf("%s│%s Worked:   %s%s%s / %s\n", Yellow, Reset, Yellow, formatDuration(worked), Reset, formatDuration(workdayTarget))
	fmt.Printf("%s│%s Remaining: %s%s%s\n", Yellow, Reset, Yellow, formatDuration(remaining), Reset)
	fmt.Printf("%s│%s Progress:  %s%d%%%s\n", Yellow, Reset, Yellow, percentWorked, Reset)
	fmt.Printf("%s└─────────────────────────────────────┘%s\n", Yellow, Reset)
	fmt.Println()

	// Next action (highlighted)
	fmt.Printf("%s╔════════════════════════════════════════╗%s\n", Magenta, Reset)
	fmt.Printf("%s║%s [!] WORKLOG [!]%s%s%s\n", Magenta, Bold, Reset, strings.Repeat(" ", 20), Magenta)
	fmt.Printf("%s║%s [!] %-28s [!] %s%s\n", Magenta, Bold, fmt.Sprintf("%s - Finish by: %s", nextAction(s.Step), checkoutTime.Format("3:04 PM")), Reset, Magenta)
	fmt.Printf("%s╚════════════════════════════════════════╝%s\n", Magenta, Reset)
	fmt.Println()
}

func printBanner() {
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("WORKLOG", pterm.NewStyle(pterm.FgRed, pterm.Bold)),
	).Render()
	fmt.Printf("%s[!] Penetrate your schedule%s\n", Red, Reset)
	fmt.Printf("%s[!] Status: Engaged%s\n", Red, Reset)
}

func getStatusSymbol(step int) string {
	switch step {
	case 0:
		return "⭕"
	case 1:
		return "🟢"
	case 2:
		return "🟡"
	case 3:
		return "🟢"
	case 4:
		return "✓"
	default:
		return "?"
	}
}

func getStatusColor(step int) string {
	switch step {
	case 0:
		return Red
	case 1:
		return Green
	case 2:
		return Yellow
	case 3:
		return Green
	case 4:
		return Cyan
	default:
		return Reset
	}
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
	return t.Format("3:04 PM")
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
