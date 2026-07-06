package main

import (
	"fmt"
	"os"
	"time"

	"worklog/internal/actions"
	"worklog/internal/state"
	"worklog/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: worklog chkin [HH:MM] | chkout [HH:MM] | stat | reset | update")
		return
	}

	cmd := os.Args[1]
	s, err := state.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to load state: %v\n", err)
	}

	arg := ""
	if len(os.Args) > 2 {
		arg = os.Args[2]
	}

	var ref time.Time
	if cmd == "chkin" || cmd == "checkin" {
		if s.Step == 2 {
			ref, _ = time.Parse(time.RFC3339, s.LunchStart)
		}
	} else if cmd == "chkout" || cmd == "checkout" {
		if s.Step == 1 {
			ref, _ = time.Parse(time.RFC3339, s.EntryTime)
		} else if s.Step == 3 {
			ref, _ = time.Parse(time.RFC3339, s.LunchEnd)
		}
	}

	switch cmd {
	case "chkin", "checkin":
		ts, err := actions.ParseOptionalTime(arg, ref)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid time format, use HH:MM")
			return
		}
		if err := actions.Chkin(&s, ts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if err := state.Save(s); err != nil {
			fmt.Fprintln(os.Stderr, "failed to save state:", err)
			os.Exit(1)
		}

	case "chkout", "checkout":
		ts, err := actions.ParseOptionalTime(arg, ref)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid time format, use HH:MM")
			return
		}
		if err := actions.Chkout(&s, ts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if err := state.Save(s); err != nil {
			fmt.Fprintln(os.Stderr, "failed to save state:", err)
			os.Exit(1)
		}

	case "stat", "status":
		ui.Progress(s)
		fmt.Println("\n--- LOG ---")
		for _, l := range s.Logs {
			fmt.Println(l)
		}

	case "reset":
		if err := state.Reset(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to reset state:", err)
			os.Exit(1)
		}
		fmt.Println("state reset")

	case "update":
		if err := actions.Update(); err != nil {
			fmt.Fprintln(os.Stderr, "update failed:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("unknown command")
	}
}
