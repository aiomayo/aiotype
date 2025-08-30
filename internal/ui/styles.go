package ui

import "github.com/charmbracelet/lipgloss"

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
