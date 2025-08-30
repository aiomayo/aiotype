package results

import (
	"fmt"

	"aiotype/internal/ui/shared"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.result == nil {
		return ""
	}

	title := shared.ResultTitleStyle.Render("🎉 Test Complete!")

	stats := []string{
		fmt.Sprintf("%s %s", shared.StatLabelStyle.Render("WPM:"), shared.StatValueStyle.Render(fmt.Sprintf("%.1f", m.result.WPM))),
		fmt.Sprintf("%s %s", shared.StatLabelStyle.Render("Accuracy:"), shared.StatValueStyle.Render(fmt.Sprintf("%.1f%%", m.result.Accuracy))),
		fmt.Sprintf("%s %s", shared.StatLabelStyle.Render("Words:"), shared.StatValueStyle.Render(fmt.Sprintf("%d/%d", m.result.CorrectWords, m.result.TotalWords))),
		fmt.Sprintf("%s %s", shared.StatLabelStyle.Render("Characters:"), shared.StatValueStyle.Render(fmt.Sprintf("%d/%d", m.result.CorrectChars, m.result.TotalChars))),
		fmt.Sprintf("%s %s", shared.StatLabelStyle.Render("Errors:"), shared.StatValueStyle.Render(fmt.Sprintf("%d", m.result.ErrorCount))),
	}

	statsDisplay := lipgloss.JoinVertical(lipgloss.Left, stats...)
	help := shared.HelpStyle.Render("ENTER/R to restart • ESC for menu • Q to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		statsDisplay,
		"",
		help,
	)

	container := shared.ResultsContainerStyle.Width(50).Render(content)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}
