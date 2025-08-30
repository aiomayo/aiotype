package main

import (
	"fmt"
	"os"

	"aiotype/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("aiotype %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s by %s\n", date, builtBy)
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("aiotype - monkeytype but in Terminal")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  aiotype           Start the typing test")
		fmt.Println("  aiotype --version Show version information")
		fmt.Println("  aiotype --help    Show this help message")
		fmt.Println()
		return
	}

	model := ui.NewModel()

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}