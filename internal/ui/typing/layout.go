package typing

import (
	"strings"

	"aiotype/internal"
	"aiotype/internal/ui/typing/components"
)

type LayoutConstraints struct {
	minTextWidth  int
	maxTextWidth  int
	textAreaWidth int
	contentWidth  int
	boxWidth      int
	paddingTotal  int
	safetyMargin  int
}

func (m *Model) calculateOptimalLayout() LayoutConstraints {
	terminalWidth := m.windowWidth
	terminalHeight := m.windowHeight

	if terminalWidth < MinTerminalWidth || terminalHeight < MinTerminalHeight {
		return LayoutConstraints{
			minTextWidth:  MinTextAreaWidth,
			maxTextWidth:  MinTextAreaWidth,
			textAreaWidth: MinTextAreaWidth,
			contentWidth:  MinTextAreaWidth + TextPadding,
			boxWidth:      MinTextAreaWidth + TextPadding + BorderWidth,
			paddingTotal:  BorderWidth + ContentPadding + SidePadding,
			safetyMargin:  MinSafetyPadding,
		}
	}

	totalReserved := BorderWidth + TextPadding + SidePadding
	availableWidth := terminalWidth - totalReserved

	textAreaWidth := availableWidth
	if textAreaWidth < MinTextAreaWidth {
		textAreaWidth = MinTextAreaWidth
	}
	if textAreaWidth > MaxTextAreaWidth {
		textAreaWidth = MaxTextAreaWidth
	}

	contentWidth := textAreaWidth + TextPadding
	boxWidth := contentWidth + BorderWidth

	return LayoutConstraints{
		minTextWidth:  MinTextAreaWidth,
		maxTextWidth:  MaxTextAreaWidth,
		textAreaWidth: textAreaWidth,
		contentWidth:  contentWidth,
		boxWidth:      boxWidth,
		paddingTotal:  totalReserved,
		safetyMargin:  SidePadding,
	}
}

func (m *Model) processContentForBorder(rawContent string, constraints LayoutConstraints) []string {
	if rawContent == "" {
		return []string{strings.Repeat(" ", constraints.contentWidth)}
	}

	rawLines := strings.Split(rawContent, "\n")
	processedLines := make([]string, 0, len(rawLines)+ExtraLinesPadding)

	processedLines = append(processedLines, strings.Repeat(" ", constraints.contentWidth))

	for _, line := range rawLines {
		paddedLine := strings.Repeat(" ", TextPadding) + line
		processedLines = append(processedLines, paddedLine)
	}

	processedLines = append(processedLines, strings.Repeat(" ", constraints.contentWidth))
	return processedLines
}

func (m *Model) renderTypingBox(processedContent []string, timeLeft string, wpmText string, constraints LayoutConstraints) string {
	if len(processedContent) == 0 {
		return "No content"
	}

	boxWidth := constraints.boxWidth
	boxHeight := len(processedContent) + BorderWidth
	perimeter := 2*boxWidth + 2*boxHeight - PerimeterBorderAdjust
	correctEnd, errorEnd := m.calculateProgressBounds(perimeter)

	pb := components.ProgressBorder{
		Width:        boxWidth,
		Height:       boxHeight,
		ContentWidth: constraints.contentWidth,
		Content:      processedContent,
		Perimeter:    perimeter,
		CorrectEnd:   correctEnd,
		ErrorEnd:     errorEnd,
		TimeLeft:     timeLeft,
		WpmText:      wpmText,
	}

	return pb.Render()
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

func (m *Model) wrapText(text string, width int) []string {
	if width < MinTextWidth {
		width = MinTextWidth
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
		for i := lineEnd - 1; i > lineStart && i > lineEnd-WordBreakLookahead; i-- {
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
