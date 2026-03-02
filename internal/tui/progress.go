package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ProgressModel struct {
	lines []string
	done  bool
	err   error
}

type ProgressMsg string
type DoneMsg error

func NewProgressModel() ProgressModel {
	return ProgressModel{
		lines: []string{},
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		m.lines = append(m.lines, string(msg))
		return m, nil
	case DoneMsg:
		m.done = true
		m.err = error(msg)
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProgressModel) View() string {
	var s strings.Builder
	s.WriteString(TitleStyle.Render("◆ Basalt  v0.1.0"))
	s.WriteString("\n")
	s.WriteString(SeparatorStyle.Render(strings.Repeat("─", 45)))
	s.WriteString("\n\n")

	for _, line := range m.lines {
		s.WriteString(line + "\n")
	}

	if m.done {
		s.WriteString("\n")
		if m.err != nil {
			s.WriteString(ErrorStyle.Render(" ✗ Generation failed: " + m.err.Error()))
		} else {
			s.WriteString(SeparatorStyle.Render(strings.Repeat("─", 45)))
			s.WriteString("\n")
			s.WriteString(SuccessStyle.Render(" ✓ Done") + "  →  cd generated && npm install && npm run dev")
		}
		s.WriteString("\n")
	}

	return s.String()
}
