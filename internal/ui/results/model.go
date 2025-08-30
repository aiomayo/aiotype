package results

import (
	"aiotype/internal"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	result       *internal.TestResult
	windowWidth  int
	windowHeight int
}

func NewModel(result *internal.TestResult) *Model {
	return &Model{
		result: result,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	}
	return m, nil
}

func (m *Model) SetResult(result *internal.TestResult) {
	m.result = result
}
