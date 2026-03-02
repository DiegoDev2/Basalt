package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	ErrorColor     = lipgloss.Color("#FF4C4C")
	GrayColor      = lipgloss.Color("#626262")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginLeft(1).
			MarginTop(1).
			MarginBottom(1)

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(GrayColor).
			MarginBottom(1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff"))

	UrlStyle = lipgloss.NewStyle().
			Foreground(GrayColor).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(GrayColor).
			Width(12)

	PrimaryStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	GrayStyle = lipgloss.NewStyle().
			Foreground(GrayColor)
)
