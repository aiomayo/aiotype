package components

import (
	"fmt"
	"strings"

	"aiotype/internal/ui/shared"
	"github.com/charmbracelet/lipgloss"
)

type ProgressBorder struct {
	Width        int
	Height       int
	ContentWidth int
	Content      []string

	Perimeter  int
	CorrectEnd int
	ErrorEnd   int

	TimeLeft string
	WpmText  string
}

func (pb *ProgressBorder) GetBorderStyleForPosition(pos int) lipgloss.Style {
	if pos < 0 || pos >= pb.Perimeter {
		return shared.GrayTextStyle
	}

	if pos < pb.CorrectEnd {
		return shared.WhiteTextStyle
	}
	if pos < pb.ErrorEnd {
		return shared.RedTextStyle
	}
	return shared.GrayTextStyle
}

func (pb *ProgressBorder) Render() string {
	if pb.Width < 5 || pb.Height < 3 {
		return "Terminal too small"
	}

	result := strings.Builder{}
	result.Grow(pb.Width * pb.Height * 2)

	for y := 0; y < pb.Height; y++ {
		if y > 0 {
			result.WriteString("\n")
		}

		if y == 0 {
			result.WriteString(pb.renderTopBorderWithStats())
		} else if y == pb.Height-1 {
			result.WriteString(pb.renderBottomBorder())
		} else {
			result.WriteString(pb.renderContentLine(y))
		}
	}

	return result.String()
}

func (pb *ProgressBorder) renderTopBorderWithStats() string {
	statsText := pb.formatStatsText()
	availableSpace := pb.Width - 6

	if len(statsText) > availableSpace {
		return pb.renderSimpleBorder("╭", "─", "╮")
	}

	return pb.buildTopBorderWithStats(statsText)
}

func (pb *ProgressBorder) formatStatsText() string {
	statsText := fmt.Sprintf("Time: %s | WPM: %s", pb.TimeLeft, pb.WpmText)
	availableSpace := pb.Width - 6

	if len(statsText) > availableSpace {
		return fmt.Sprintf("%s|%s", pb.TimeLeft, pb.WpmText)
	}
	return statsText
}

func (pb *ProgressBorder) buildTopBorderWithStats(statsText string) string {
	nonDashChars := 5 + len(statsText)
	remainingDashes := pb.Width - nonDashChars

	if remainingDashes < 0 {
		return pb.renderSimpleBorder("╭", "─", "╮")
	}

	var result strings.Builder

	result.WriteString(pb.GetBorderStyleForPosition(0).Render("╭"))
	result.WriteString(pb.GetBorderStyleForPosition(1).Render("─"))
	result.WriteString(" ")
	result.WriteString(statsText)
	result.WriteString(" ")

	for i := 0; i < remainingDashes; i++ {
		pos := 1 + 1 + len(statsText) + 1 + i
		result.WriteString(pb.GetBorderStyleForPosition(pos).Render("─"))
	}

	result.WriteString(pb.GetBorderStyleForPosition(pb.Width - 1).Render("╮"))

	return result.String()
}

func (pb *ProgressBorder) renderSimpleBorder(left, middle, right string) string {
	var result strings.Builder
	result.WriteString(pb.GetBorderStyleForPosition(0).Render(left))

	for i := 1; i < pb.Width-1; i++ {
		result.WriteString(pb.GetBorderStyleForPosition(i).Render(middle))
	}

	result.WriteString(pb.GetBorderStyleForPosition(pb.Width - 1).Render(right))
	return result.String()
}

func (pb *ProgressBorder) renderBottomBorder() string {
	var result strings.Builder

	for x := 0; x < pb.Width; x++ {
		var pos int
		if x == 0 {
			pos = pb.Width + (pb.Height - 2) + (pb.Width - 1)
			result.WriteString(pb.GetBorderStyleForPosition(pos).Render("╰"))
		} else if x == pb.Width-1 {
			pos = pb.Width + (pb.Height - 2)
			result.WriteString(pb.GetBorderStyleForPosition(pos).Render("╯"))
		} else {
			pos = pb.Width + (pb.Height - 2) + (pb.Width - 1 - x)
			result.WriteString(pb.GetBorderStyleForPosition(pos).Render("─"))
		}
	}

	return result.String()
}

func (pb *ProgressBorder) renderContentLine(y int) string {
	var result strings.Builder

	leftPos := 2*pb.Width + pb.Height - 2 + (pb.Height - 1 - y)
	result.WriteString(pb.GetBorderStyleForPosition(leftPos).Render("│"))

	contentY := y - 1
	if contentY < len(pb.Content) {
		result.WriteString(EnforceExactWidth(pb.Content[contentY], pb.ContentWidth))
	} else {
		result.WriteString(strings.Repeat(" ", pb.ContentWidth))
	}

	rightPos := pb.Width + (y - 1)
	result.WriteString(pb.GetBorderStyleForPosition(rightPos).Render("│"))

	return result.String()
}

func EnforceExactWidth(content string, targetWidth int) string {
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
