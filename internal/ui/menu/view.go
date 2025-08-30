package menu

import (
	"aiotype/internal/ui/shared"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	title := shared.TitleStyle.Render("aiotype")
	subtitle := shared.SubtitleStyle.Render("A Typoing game inspired by monkeytype")
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669")).
		Align(lipgloss.Center).
		Render("Press ENTER or SPACE to start typing • Q to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		instructions,
	)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
