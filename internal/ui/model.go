package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aiotype/internal"
)

type Model struct {
	state         internal.GameState
	currentTest   *internal.TypingTest
	result        *internal.TestResult
	config        internal.GameConfig
	windowWidth   int
	windowHeight  int
	realTimeWPM   float64
	fadeStartTime time.Time
	mu            sync.RWMutex
}

type TickMsg time.Time

func hexToRGB(hex string) (int, int, int) {
	hex = hex[1:]
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return int(r), int(g), int(b)
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func interpolateColor(color1, color2 string, factor float64) string {
	if factor <= 0 {
		return color1
	}
	if factor >= 1 {
		return color2
	}
	r1, g1, b1 := hexToRGB(color1)
	r2, g2, b2 := hexToRGB(color2)
	r := int(float64(r1) + factor*float64(r2-r1))
	g := int(float64(g1) + factor*float64(g2-g1))
	b := int(float64(b1) + factor*float64(b2-b1))
	return rgbToHex(r, g, b)
}

func (m *Model) getFadeFactor() float64 {
	if !m.isInLastFiveSeconds() {
		return 0.0
	}

	elapsed := time.Since(m.fadeStartTime)
	phase := elapsed.Seconds() / (1500 * time.Millisecond).Seconds() * 2 * math.Pi
	return (math.Sin(phase) + 1.0) / 2.0
}

func (m *Model) getCurrentCursorStyle() lipgloss.Style {
	if !m.isInLastFiveSeconds() {
		return CursorStyle
	}
	fadeFactor := m.getFadeFactor()
	interpolatedColor := interpolateColor("#ffffff", "#ff0000", fadeFactor)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color(interpolatedColor))
}

func tickEvery() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func NewModel() *Model {
	config := internal.GameConfig{
		TestDuration: 10 * time.Second,
		WordCount:    50,
	}

	return &Model{
		state:  internal.StateMenu,
		config: config,
	}
}

func (m *Model) Init() tea.Cmd {
	return tickEvery()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, tea.ClearScreen

	case tea.KeyMsg:
		return m.handleKeyInput(msg)

	case TickMsg:
		m.mu.RLock()
		if m.state == internal.StateTyping && m.currentTest != nil && !m.currentTest.StartTime.IsZero() {
			wpm := internal.CalculateWPM(m.currentTest)
			m.mu.RUnlock()
			m.mu.Lock()
			m.realTimeWPM = wpm

			if m.isInLastFiveSeconds() {
				if m.fadeStartTime.IsZero() {
					m.fadeStartTime = time.Now()
				}
			} else {
				m.fadeStartTime = time.Time{}
			}

			m.mu.Unlock()
		} else {
			m.mu.RUnlock()
		}
		return m, tickEvery()
	}

	return m, nil
}

func (m *Model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case internal.StateMenu:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			m.startNewTest()
			return m, nil
		}

	case internal.StateTyping:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mu.Lock()
			m.state = internal.StateMenu
			m.currentTest = nil
			m.mu.Unlock()
			return m, nil
		case "backspace":
			m.mu.Lock()
			if m.currentTest != nil {
				internal.ProcessBackspace(m.currentTest)
			}
			m.mu.Unlock()
			return m, nil
		default:
			m.mu.Lock()
			if m.currentTest != nil {
				msgStr := msg.String()
				runes := []rune(msgStr)
				if len(runes) == 1 {
					char := runes[0]
					completed := internal.ProcessCharacter(m.currentTest, char)
					if completed {
						m.result = internal.GenerateResult(m.currentTest)
						m.state = internal.StateResults
					}
				}
			}
			m.mu.Unlock()
			return m, nil
		}

	case internal.StateResults:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ", "r":
			m.startNewTest()
			return m, nil
		case "esc":
			m.state = internal.StateMenu
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) startNewTest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTest = internal.NewTest(m.config)
	m.state = internal.StateTyping
	m.result = nil
	m.realTimeWPM = 0
	m.fadeStartTime = time.Time{}
}

func (m *Model) View() string {
	switch m.state {
	case internal.StateMenu:
		return m.renderMenu()
	case internal.StateTyping:
		return m.renderTyping()
	case internal.StateResults:
		return m.renderResults()
	}
	return ""
}

func (m *Model) renderMenu() string {
	title := TitleStyle.Render("🐵 AIOType")
	subtitle := SubtitleStyle.Render("A beautiful typing test inspired by Monkey Type")
	instructions := lipgloss.NewStyle().
		Foreground(grayColor).
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

func (m *Model) renderTyping() string {
	if m.currentTest == nil {
		return ""
	}

	typingArea := m.renderTypingArea()
	help := HelpStyle.Render("ESC to return to menu • CTRL+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		typingArea,
		help,
	)

	styledContent := BaseStyle.Render(content)

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
	wrappedLines := m.wrapText(m.currentTest.TargetText, constraints.maxWidth-4)
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

func (m *Model) isInLastFiveSeconds() bool {
	if m.currentTest == nil || m.currentTest.StartTime.IsZero() {
		return false
	}

	elapsed := time.Since(m.currentTest.StartTime)
	remaining := m.config.TestDuration - elapsed
	return remaining <= 5*time.Second && remaining > 0
}

func (m *Model) processContentForBorder(rawContent string, constraints LayoutConstraints) []string {
	if rawContent == "" {
		return []string{strings.Repeat(" ", constraints.maxWidth)}
	}

	rawLines := strings.Split(rawContent, "\n")
	processedLines := make([]string, 0, len(rawLines)+4)

	processedLines = append(processedLines, strings.Repeat(" ", constraints.maxWidth))

	for _, line := range rawLines {
		paddedLine := "  " + line
		processedLine := enforceExactWidth(paddedLine, constraints.maxWidth)
		processedLines = append(processedLines, processedLine)
	}

	processedLines = append(processedLines, strings.Repeat(" ", constraints.maxWidth))
	return processedLines
}

func (m *Model) renderTypingBox(processedContent []string, timeLeft string, wpmText string, constraints LayoutConstraints) string {
	if len(processedContent) == 0 {
		return "No content"
	}

	boxWidth := constraints.maxWidth + 2
	boxHeight := len(processedContent) + 2
	perimeter := 2*boxWidth + 2*boxHeight - 4
	correctEnd, errorEnd := m.calculateProgressBounds(perimeter)

	pb := ProgressBorder{
		width:        boxWidth,
		height:       boxHeight,
		contentWidth: constraints.maxWidth,
		content:      processedContent,
		perimeter:    perimeter,
		correctEnd:   correctEnd,
		errorEnd:     errorEnd,
		timeLeft:     timeLeft,
		wpmText:      wpmText,
	}

	return pb.render()
}

func (m *Model) calculateProgressBounds(perimeter int) (correctEnd, errorEnd int) {
	if m.currentTest == nil {
		return 0, 0
	}

	correctChars := internal.CountCorrectChars(m.currentTest)
	totalChars := len(m.currentTest.TargetText)

	if totalChars == 0 {
		return 0, 0
	}

	correctProgress := float64(correctChars) / float64(totalChars)
	totalProgress := float64(m.currentTest.CurrentPos) / float64(totalChars)

	correctEnd = int(correctProgress * float64(perimeter))
	errorEnd = int(totalProgress * float64(perimeter))

	if correctEnd < 0 {
		correctEnd = 0
	}
	if correctEnd > perimeter {
		correctEnd = perimeter
	}
	if errorEnd < correctEnd {
		errorEnd = correctEnd
	}
	if errorEnd > perimeter {
		errorEnd = perimeter
	}

	return correctEnd, errorEnd
}

func enforceExactWidth(content string, targetWidth int) string {
	if targetWidth <= 0 {
		return ""
	}

	displayWidth := lipgloss.Width(content)

	if displayWidth >= targetWidth {
		return content
	}

	padding := targetWidth - displayWidth
	return content + strings.Repeat(" ", padding)
}

type LayoutConstraints struct {
	minWidth     int
	maxWidth     int
	paddingTotal int
	safetyMargin int
}

func (m *Model) calculateOptimalLayout() LayoutConstraints {
	terminalWidth := m.windowWidth

	totalReserved := 2 + 4 + 4
	contentWidth := terminalWidth - totalReserved

	if contentWidth < 20 {
		contentWidth = 20
	}

	return LayoutConstraints{
		minWidth:     20,
		maxWidth:     contentWidth,
		paddingTotal: totalReserved,
		safetyMargin: 4,
	}
}

type ProgressBorder struct {
	width        int
	height       int
	contentWidth int
	content      []string

	perimeter  int
	correctEnd int
	errorEnd   int

	timeLeft string
	wpmText  string
}

func (pb *ProgressBorder) getBorderStyleForPosition(pos int) lipgloss.Style {
	if pos < 0 || pos >= pb.perimeter {
		return lipgloss.NewStyle().Foreground(grayColor)
	}

	if pos < pb.correctEnd {
		return lipgloss.NewStyle().Foreground(whiteColor)
	}
	if pos < pb.errorEnd {
		return lipgloss.NewStyle().Foreground(redColor)
	}
	return lipgloss.NewStyle().Foreground(grayColor)
}

func (pb *ProgressBorder) render() string {
	if pb.width < 5 || pb.height < 3 {
		return "Terminal too small"
	}

	result := strings.Builder{}
	result.Grow(pb.width * pb.height * 2)

	for y := 0; y < pb.height; y++ {
		if y > 0 {
			result.WriteString("\n")
		}

		if y == 0 {
			result.WriteString(pb.renderTopBorderWithStats())
		} else if y == pb.height-1 {
			result.WriteString(pb.renderBottomBorder())
		} else {
			result.WriteString(pb.renderContentLine(y))
		}
	}

	return result.String()
}

func (pb *ProgressBorder) renderTopBorderWithStats() string {
	statsText := pb.formatStatsText()
	availableSpace := pb.width - 6

	if len(statsText) > availableSpace {
		return pb.renderSimpleBorder("╭", "─", "╮")
	}

	return pb.buildTopBorderWithStats(statsText)
}

func (pb *ProgressBorder) formatStatsText() string {
	statsText := fmt.Sprintf("Time: %s | WPM: %s", pb.timeLeft, pb.wpmText)
	availableSpace := pb.width - 6

	if len(statsText) > availableSpace {
		return fmt.Sprintf("%s|%s", pb.timeLeft, pb.wpmText)
	}
	return statsText
}

func (pb *ProgressBorder) buildTopBorderWithStats(statsText string) string {
	nonDashChars := 5 + len(statsText)
	remainingDashes := pb.width - nonDashChars

	if remainingDashes < 0 {
		return pb.renderSimpleBorder("╭", "─", "╮")
	}

	var result strings.Builder

	result.WriteString(pb.getBorderStyleForPosition(0).Render("╭"))
	result.WriteString(pb.getBorderStyleForPosition(1).Render("─"))
	result.WriteString(" ")
	result.WriteString(statsText)
	result.WriteString(" ")

	for i := 0; i < remainingDashes; i++ {
		pos := 1 + 1 + len(statsText) + 1 + i
		result.WriteString(pb.getBorderStyleForPosition(pos).Render("─"))
	}

	result.WriteString(pb.getBorderStyleForPosition(pb.width - 1).Render("╮"))

	return result.String()
}

func (pb *ProgressBorder) renderSimpleBorder(left, middle, right string) string {
	var result strings.Builder
	result.WriteString(pb.getBorderStyleForPosition(0).Render(left))

	for i := 1; i < pb.width-1; i++ {
		result.WriteString(pb.getBorderStyleForPosition(i).Render(middle))
	}

	result.WriteString(pb.getBorderStyleForPosition(pb.width - 1).Render(right))
	return result.String()
}

func (pb *ProgressBorder) renderBottomBorder() string {
	var result strings.Builder

	for x := 0; x < pb.width; x++ {
		var pos int
		if x == 0 {
			pos = pb.width + (pb.height - 2) + (pb.width - 1)
			result.WriteString(pb.getBorderStyleForPosition(pos).Render("╰"))
		} else if x == pb.width-1 {
			pos = pb.width + (pb.height - 2)
			result.WriteString(pb.getBorderStyleForPosition(pos).Render("╯"))
		} else {
			pos = pb.width + (pb.height - 2) + (pb.width - 1 - x)
			result.WriteString(pb.getBorderStyleForPosition(pos).Render("─"))
		}
	}

	return result.String()
}

func (pb *ProgressBorder) renderContentLine(y int) string {
	var result strings.Builder

	leftPos := 2*pb.width + pb.height - 2 + (pb.height - 1 - y)
	result.WriteString(pb.getBorderStyleForPosition(leftPos).Render("│"))

	contentY := y - 1
	if contentY < len(pb.content) {
		result.WriteString(enforceExactWidth(pb.content[contentY], pb.contentWidth))
	} else {
		result.WriteString(strings.Repeat(" ", pb.contentWidth))
	}

	rightPos := pb.width + (y - 1)
	result.WriteString(pb.getBorderStyleForPosition(rightPos).Render("│"))

	return result.String()
}

func (m *Model) wrapText(text string, width int) []string {
	if width < 10 {
		width = 10
	}

	if len(text) == 0 {
		return []string{""}
	}

	var lines []string
	textRunes := []rune(text)
	lineStart := 0

	for lineStart < len(textRunes) {
		lineEnd := lineStart + width

		if lineEnd >= len(textRunes) {
			lines = append(lines, string(textRunes[lineStart:]))
			break
		}

		breakPoint := lineEnd
		for i := lineEnd - 1; i > lineStart && i > lineEnd-20; i-- {
			if textRunes[i] == ' ' {
				breakPoint = i + 1
				break
			}
		}

		lines = append(lines, string(textRunes[lineStart:breakPoint]))
		lineStart = breakPoint
	}

	if len(lines) == 0 {
		lines = []string{""}
	}

	return lines
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

		for _, char := range line {
			if charIndex < len(m.currentTest.TypedChars) {
				isInErrorUnit := false
				for _, unitStatus := range m.currentTest.WordStatuses {
					if unitStatus.IsComplete && unitStatus.HasError {
						if charIndex >= unitStatus.StartIndex && charIndex <= unitStatus.EndIndex {
							isInErrorUnit = true
							break
						}
					}
				}

				typedChar := m.currentTest.TypedChars[charIndex]
				if isInErrorUnit {
					if typedChar.IsCorrect {
						result.WriteString(WhiteTextRedBgStyle.Render(string(typedChar.Character)))
					} else {
						result.WriteString(RedTextRedBgStyle.Render(string(typedChar.Character)))
					}
				} else if typedChar.IsCorrect {
					result.WriteString(WhiteTextStyle.Render(string(typedChar.Character)))
				} else {
					result.WriteString(RedTextStyle.Render(string(typedChar.Character)))
				}
			} else if charIndex == m.currentTest.CurrentPos {
				cursorStyle := m.getCurrentCursorStyle()
				result.WriteString(cursorStyle.Render(string(char)))
			} else {
				result.WriteString(GrayTextStyle.Render(string(char)))
			}
			charIndex++
		}
	}

	return result.String()
}

func (m *Model) renderResults() string {
	if m.result == nil {
		return ""
	}

	title := ResultTitleStyle.Render("🎉 Test Complete!")

	stats := []string{
		fmt.Sprintf("%s %s", StatLabelStyle.Render("WPM:"), StatValueStyle.Render(fmt.Sprintf("%.1f", m.result.WPM))),
		fmt.Sprintf("%s %s", StatLabelStyle.Render("Accuracy:"), StatValueStyle.Render(fmt.Sprintf("%.1f%%", m.result.Accuracy))),
		fmt.Sprintf("%s %s", StatLabelStyle.Render("Words:"), StatValueStyle.Render(fmt.Sprintf("%d/%d", m.result.CorrectWords, m.result.TotalWords))),
		fmt.Sprintf("%s %s", StatLabelStyle.Render("Characters:"), StatValueStyle.Render(fmt.Sprintf("%d/%d", m.result.CorrectChars, m.result.TotalChars))),
		fmt.Sprintf("%s %s", StatLabelStyle.Render("Errors:"), StatValueStyle.Render(fmt.Sprintf("%d", m.result.ErrorCount))),
	}

	statsDisplay := lipgloss.JoinVertical(lipgloss.Left, stats...)
	help := HelpStyle.Render("ENTER/R to restart • ESC for menu • Q to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		statsDisplay,
		"",
		help,
	)

	container := ResultsContainerStyle.Width(50).Render(content)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}
