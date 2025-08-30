package shared

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	colorGray    = "#646669"
	colorWhite   = "#ffffff"
	colorRed     = "#ff0000"
	colorPrimary = "#e2b714"
	colorBlack   = "#000000"
	colorDarkRed = "#4d0000"
)

var (
	grayColor    = lipgloss.Color(colorGray)
	whiteColor   = lipgloss.Color(colorWhite)
	redColor     = lipgloss.Color(colorRed)
	primaryColor = lipgloss.Color(colorPrimary)
	darkRedColor = lipgloss.Color(colorDarkRed)
)

var (
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 0).
			Foreground(whiteColor)

	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Align(lipgloss.Center).
			MarginBottom(2)

	GrayTextStyle = lipgloss.NewStyle().
			Foreground(grayColor)

	WhiteTextStyle = lipgloss.NewStyle().
			Foreground(whiteColor)

	RedTextStyle = lipgloss.NewStyle().
			Foreground(redColor)

	WhiteTextRedBgStyle = lipgloss.NewStyle().
				Foreground(whiteColor).
				Background(darkRedColor)

	RedTextRedBgStyle = lipgloss.NewStyle().
				Foreground(redColor).
				Background(darkRedColor)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorBlack)).
			Background(whiteColor)

	StatLabelStyle = lipgloss.NewStyle().
			Foreground(grayColor)

	StatValueStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	ResultsContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginTop(2)

	ResultTitleStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Align(lipgloss.Center).
				MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Align(lipgloss.Center).
			MarginTop(1)
)

func HexToRGB(hex string) (int, int, int) {
	hex = hex[1:]
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return int(r), int(g), int(b)
}

func RgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func InterpolateColor(color1, color2 string, factor float64) string {
	if factor <= 0 {
		return color1
	}
	if factor >= 1 {
		return color2
	}
	r1, g1, b1 := HexToRGB(color1)
	r2, g2, b2 := HexToRGB(color2)
	r := int(float64(r1) + factor*float64(r2-r1))
	g := int(float64(g1) + factor*float64(g2-g1))
	b := int(float64(b1) + factor*float64(b2-b1))
	return RgbToHex(r, g, b)
}
