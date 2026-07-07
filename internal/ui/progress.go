package ui

import (
	"fmt"
	"os"
	"os/signal"
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
	fmt.Println()
	printBanner()
	fmt.Println()

	worked := s.WorkedDuration()
	remaining := workdayTarget - worked
	if remaining < 0 {
		remaining = 0
	}
	checkoutTime := time.Now().Add(remaining)
	percentWorked := int(percent(worked, workdayTarget))

	fmt.Print(buildProgressString(s, checkoutTime, worked, remaining, percentWorked))
	fmt.Println()
}

func LiveProgress(s state.State) {
	clearScreen()
	printBanner()
	fmt.Println()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	defer signal.Stop(stop)

	area, _ := pterm.DefaultArea.Start()
	defer area.Stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		worked := s.WorkedDuration()
		remaining := workdayTarget - worked
		if remaining < 0 {
			remaining = 0
		}
		checkoutTime := time.Now().Add(remaining)
		percentWorked := int(percent(worked, workdayTarget))

		content := buildProgressString(s, checkoutTime, worked, remaining, percentWorked)
		content += fmt.Sprintf("\n%s[Live] Press Ctrl+C to stop live status refresh.%s\n", Cyan, Reset)

		area.Update(content)

		select {
		case <-ticker.C:
			// Continue loop
		case <-stop:
			return
		}
	}
}

func buildProgressString(s state.State, checkoutTime time.Time, worked, remaining time.Duration, percentWorked int) string {
	var sb strings.Builder

	// Date and Status
	sb.WriteString(fmt.Sprintf("%s▸ DATE:%s      %s\n", Cyan, Reset, s.Date))
	statusSymbol := getStatusSymbol(s.Step)
	statusColor := getStatusColor(s.Step)
	sb.WriteString(fmt.Sprintf("%s▸ STATUS:%s    %s%s%s %s%s\n\n", Cyan, Reset, statusColor, statusSymbol, Reset, Bold, label(s.Step)))

	// Check in and out times
	checkInTime := formatTime(s.EntryTime)
	checkOutTime := formatTime(s.ExitTime)
	sb.WriteString(fmt.Sprintf("%s┌─ TIMES ─────────────────────────────┐%s\n", Green, Reset))
	sb.WriteString(fmt.Sprintf("%s│%s CHECK IN:  %s%s%s\n", Green, Reset, Green, checkInTime, Reset))
	if checkOutTime == "-" {
		sb.WriteString(fmt.Sprintf("%s│%s CHECK OUT: %s%s (est)%s\n", Green, Reset, Yellow, checkoutTime.Format("3:04 PM"), Reset))
	} else {
		sb.WriteString(fmt.Sprintf("%s│%s CHECK OUT: %s%s%s\n", Green, Reset, Green, checkOutTime, Reset))
	}
	sb.WriteString(fmt.Sprintf("%s└─────────────────────────────────────┘%s\n\n", Green, Reset))

	// Lunch times
	lunchStart := formatTime(s.LunchStart)
	lunchEnd := formatTime(s.LunchEnd)
	sb.WriteString(fmt.Sprintf("%s┌─ LUNCH ─────────────────────────────┐%s\n", Blue, Reset))
	sb.WriteString(fmt.Sprintf("%s│%s START: %s%s   END: %s%s%s\n", Blue, Reset, Blue, lunchStart, Blue, lunchEnd, Reset))
	sb.WriteString(fmt.Sprintf("%s└─────────────────────────────────────┘%s\n\n", Blue, Reset))

	// Work hours with visual indicator
	progressBar := renderProgressBar(percentWorked)
	sb.WriteString(fmt.Sprintf("%s┌─ HOURS ─────────────────────────────┐%s\n", Yellow, Reset))
	sb.WriteString(fmt.Sprintf("%s│%s Worked:    %s%s%s / %s\n", Yellow, Reset, Yellow, formatDuration(worked), Reset, formatDuration(workdayTarget)))
	sb.WriteString(fmt.Sprintf("%s│%s Remaining: %s%s%s\n", Yellow, Reset, Yellow, formatDuration(remaining), Reset))
	sb.WriteString(fmt.Sprintf("%s│%s Checkout:  %s%s%s\n", Yellow, Reset, Yellow, checkoutTime.Format("3:04 PM"), Reset))
	sb.WriteString(fmt.Sprintf("%s│%s Progress:  [%s] %s%d%%%s\n", Yellow, Reset, progressBar, Yellow, percentWorked, Reset))
	sb.WriteString(fmt.Sprintf("%s└─────────────────────────────────────┘%s\n\n", Yellow, Reset))

	// Next action
	sb.WriteString(fmt.Sprintf("%s╔════════════════════════════════════════╗%s\n", Magenta, Reset))
	sb.WriteString(fmt.Sprintf("%s║%s [!] WORKLOG [!]%s%s%s\n", Magenta, Bold, Reset, strings.Repeat(" ", 20), Magenta))
	sb.WriteString(fmt.Sprintf("%s║%s [!] %-28s [!] %s%s\n", Magenta, Bold, fmt.Sprintf("%s - Finish by: %s", nextAction(s.Step), checkoutTime.Format("3:04 PM")), Reset, Magenta))
	sb.WriteString(fmt.Sprintf("%s╚════════════════════════════════════════╝%s\n", Magenta, Reset))

	return sb.String()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func renderProgressBar(percent int) string {
	filled := percent / 10
	if filled > 10 {
		filled = 10
	}
	bar := strings.Repeat("=", filled) + strings.Repeat(" ", 10-filled)
	return bar
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
