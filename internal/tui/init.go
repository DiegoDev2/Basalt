package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InitStep int

const (
	StepProjectName InitStep = iota
	StepDatabase
	StepFramework
	StepAuth
	StepLanguage
	StepSummary
	StepDone
)

type InitModel struct {
	step      InitStep
	questions []question
	choices   []string
	cursor    int
	input     textinput.Model
	results   map[string]string
	quitting  bool
}

type question struct {
	name    string
	label   string
	options []option
}

type option struct {
	label string
	url   string
}

func NewInitModel() InitModel {
	ti := textinput.New()
	ti.Placeholder = "my-api"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	questions := []question{
		{
			name:  "project",
			label: "Project name",
		},
		{
			name:  "db",
			label: "Database",
			options: []option{
				{label: "Supabase", url: "https://supabase.com"},
				{label: "PostgreSQL", url: "https://postgresql.org"},
				{label: "MySQL", url: "https://mysql.com"},
			},
		},
		{
			name:  "framework",
			label: "Framework",
			options: []option{
				{label: "Express", url: "https://expressjs.com"},
				{label: "Fastify", url: "https://fastify.dev"},
				{label: "Hono", url: "https://hono.dev"},
			},
		},
		{
			name:  "auth",
			label: "Auth",
			options: []option{
				{label: "JWT", url: "https://jwt.io"},
				{label: "API Key", url: ""},
				{label: "None", url: ""},
			},
		},
		{
			name:  "lang",
			label: "Output language",
			options: []option{
				{label: "TypeScript", url: "https://typescriptlang.org"},
				{label: "JavaScript", url: ""},
			},
		},
	}

	return InitModel{
		step:      StepProjectName,
		questions: questions,
		input:     ti,
		results:   make(map[string]string),
	}
}

func (m InitModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.questions[m.step].options)-1 {
				m.cursor++
			}
		case "enter":
			if m.step == StepProjectName {
				val := m.input.Value()
				if val == "" {
					val = m.input.Placeholder
				}
				m.results["project"] = val
				m.step++
				m.cursor = 0
				return m, nil
			} else if m.step < StepSummary {
				q := m.questions[m.step]
				m.results[q.name] = q.options[m.cursor].label
				m.step++
				m.cursor = 0
				return m, nil
			} else if m.step == StepSummary {
				if m.cursor == 0 { // Yes, create project
					m.step = StepDone
					return m, tea.Quit
				} else { // No, go back
					m.step = StepProjectName
					m.cursor = 0
					return m, nil
				}
			}
		}
	}

	if m.step == StepProjectName {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m InitModel) View() string {
	if m.quitting {
		return ""
	}

	if m.step == StepDone {
		return ""
	}

	var s strings.Builder

	s.WriteString(TitleStyle.Render("◆ Basalt  v0.1.0"))
	s.WriteString("\n")
	s.WriteString(SeparatorStyle.Render(strings.Repeat("─", 45)))
	s.WriteString("\n\n")

	if m.step < StepSummary {
		// Show previous answers
		for i := 0; i < int(m.step); i++ {
			q := m.questions[i]
			s.WriteString(SuccessStyle.Render(" ✓ "))
			s.WriteString(LabelStyle.Render(q.label))
			s.WriteString(m.results[q.name])
			s.WriteString("\n")
		}

		// Show current question
		q := m.questions[m.step]
		s.WriteString(PrimaryStyle.Render(" ? "))
		s.WriteString(lipgloss.NewStyle().Bold(true).Render(q.label))
		s.WriteString("\n")

		if m.step == StepProjectName {
			s.WriteString("> ")
			s.WriteString(m.input.View())
			s.WriteString("\n")
		} else {
			for i, opt := range q.options {
				cursor := "  "
				label := opt.label
				if m.cursor == i {
					cursor = PrimaryStyle.Render(" ❯ ")
					label = SelectedStyle.Render("● " + opt.label)
				} else {
					label = UnselectedStyle.Render("○ " + opt.label)
				}
				s.WriteString(cursor)
				s.WriteString(label)
				if opt.url != "" {
					s.WriteString("          ")
					s.WriteString(UrlStyle.Render(opt.url))
				}
				s.WriteString("\n")
			}
		}
	} else if m.step == StepSummary {
		for _, q := range m.questions {
			s.WriteString(SuccessStyle.Render(" ✓ "))
			s.WriteString(LabelStyle.Render(q.label))
			s.WriteString(m.results[q.name])
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(PrimaryStyle.Render(" ? "))
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Looks good?"))
		s.WriteString("\n")

		options := []string{"Yes, create project", "No, go back"}
		for i, opt := range options {
			cursor := "  "
			if m.cursor == i {
				cursor = PrimaryStyle.Render(" ❯ ")
				s.WriteString(cursor + SelectedStyle.Render(opt) + "\n")
			} else {
				s.WriteString(cursor + UnselectedStyle.Render(opt) + "\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(SeparatorStyle.Render(strings.Repeat("─", 45)))
	s.WriteString("\n")
	s.WriteString(GrayStyle.Render(" ↑↓ navigate   enter select   ctrl+c cancel"))
	s.WriteString("\n")

	return s.String()
}

func (m InitModel) GetResults() map[string]string {
	return m.results
}

func (m InitModel) IsConfirmed() bool {
	return m.step == StepDone
}

func RunInitWizard() (map[string]string, bool, error) {
	m := NewInitModel()
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	res := finalModel.(InitModel)
	if res.quitting {
		return nil, false, nil
	}

	return res.GetResults(), res.IsConfirmed(), nil
}
