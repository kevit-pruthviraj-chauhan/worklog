package main

import (
	"fmt"
	"os"

	"worklog/internal/actions"
	"worklog/internal/state"
	"worklog/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("use: chkin | chkout | stat")
		return
	}

	cmd := os.Args[1]
	s := state.Load()

	switch cmd {
	case "chkin":
		actions.Chkin(&s)
		state.Save(s)

	case "chkout":
		actions.Chkout(&s)
		state.Save(s)

	case "stat":
		ui.Progress(s)
		fmt.Println("\n--- LOG ---")
		for _, l := range s.Logs {
			fmt.Println(l)
		}

	default:
		fmt.Println("unknown command")
	}
}