package shared

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func TickEvery() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
