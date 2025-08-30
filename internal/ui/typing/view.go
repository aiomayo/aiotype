package typing

import (
	"fmt"
	"strings"
	"time"

	"aiotype/internal/ui/shared"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.currentTest == nil {
		return ""
	}

	if m.windowWidth < MinTerminalWidth || m.windowHeight < MinTerminalHeight {
		return lipgloss.Place(
			m.windowWidth,
			m.windowHeight,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(shared.RedTextStyle.GetForeground()).Render("Terminal too small"),
		)
	}

	typingArea := m.renderTypingArea()
	help := shared.HelpStyle.Render("ESC to return to menu • CTRL+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		typingArea,
		"",
		help,
	)

	responsivePadding := BasePadding
	if m.windowWidth > LargeScreenWidth {
		responsivePadding = LargePadding
	}

	styledContent := lipgloss.NewStyle().
		Padding(responsivePadding, 0).
		Foreground(shared.WhiteTextStyle.GetForeground()).
		Render(content)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Top,
		styledContent,
	)
}

func (m *Model) renderTypingArea() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentTest == nil {
		return ""
	}

	constraints := m.calculateOptimalLayout()
	wrappedLines := m.wrapText(m.currentTest.TargetText, constraints.textAreaWidth)
	rendered := m.renderTextWithHighlighting(wrappedLines)
	processedContent := m.processContentForBorder(rendered, constraints)
	timeLeft := m.formatTimeRemaining()
	wpmText := fmt.Sprintf("%.0f", m.realTimeWPM)

	return m.renderTypingBox(processedContent, timeLeft, wpmText, constraints)
}

func (m *Model) formatTimeRemaining() string {
	if m.currentTest == nil || m.currentTest.StartTime.IsZero() {
		return fmt.Sprintf("%.0fs", m.config.TestDuration.Seconds())
	}

	elapsed := time.Since(m.currentTest.StartTime)
	remaining := m.config.TestDuration - elapsed
	if remaining <= 0 {
		return "0s"
	}
	return fmt.Sprintf("%.0fs", remaining.Seconds())
}

func (m *Model) renderTextWithHighlighting(wrappedLines []string) string {
	if m.currentTest == nil || len(wrappedLines) == 0 {
		return strings.Join(wrappedLines, "\n")
	}

	var result strings.Builder
	charIndex := 0

	for lineIndex, line := range wrappedLines {
		if lineIndex > 0 {
			result.WriteString("\n")
		}
		m.renderLineWithHighlighting(line, &result, &charIndex)
	}

	return result.String()
}

func (m *Model) renderLineWithHighlighting(line string, result *strings.Builder, charIndex *int) {
	for _, char := range line {
		m.renderCharacterWithStyle(char, result, *charIndex)
		*charIndex++
	}
}

func (m *Model) renderCharacterWithStyle(char rune, result *strings.Builder, charIndex int) {
	if charIndex < len(m.currentTest.TypedChars) {
		m.renderTypedCharacter(char, result, charIndex)
	} else if charIndex == m.currentTest.CurrentPos {
		cursorStyle := m.getCurrentCursorStyle()
		result.WriteString(cursorStyle.Render(string(char)))
	} else {
		result.WriteString(shared.GrayTextStyle.Render(string(char)))
	}
}

func (m *Model) renderTypedCharacter(char rune, result *strings.Builder, charIndex int) {
	typedChar := m.currentTest.TypedChars[charIndex]
	isInErrorUnit := m.isCharacterInErrorUnit(charIndex)

	if isInErrorUnit {
		if typedChar.IsCorrect {
			result.WriteString(shared.WhiteTextRedBgStyle.Render(string(typedChar.Character)))
		} else {
			result.WriteString(shared.RedTextRedBgStyle.Render(string(typedChar.Character)))
		}
	} else if typedChar.IsCorrect {
		result.WriteString(shared.WhiteTextStyle.Render(string(typedChar.Character)))
	} else {
		result.WriteString(shared.RedTextStyle.Render(string(typedChar.Character)))
	}
}

func (m *Model) isCharacterInErrorUnit(charIndex int) bool {
	for _, unitStatus := range m.currentTest.WordStatuses {
		if unitStatus.IsComplete && unitStatus.HasError {
			if charIndex >= unitStatus.StartIndex && charIndex <= unitStatus.EndIndex {
				return true
			}
		}
	}
	return false
}
