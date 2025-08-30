package ui

import (
	"github.com/charmbracelet/bubbletea"

	"aiotype/internal"
	"aiotype/internal/ui/menu"
	"aiotype/internal/ui/results"
	"aiotype/internal/ui/typing"
)

type Model struct {
	state        internal.GameState
	config       internal.GameConfig
	menuModel    *menu.Model
	typingModel  *typing.Model
	resultsModel *results.Model
}

func NewModel() *Model {
	config := internal.DefaultGameConfig()

	return &Model{
		state:        internal.StateMenu,
		config:       config,
		menuModel:    menu.NewModel(),
		typingModel:  typing.NewModel(config),
		resultsModel: results.NewModel(nil),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.menuModel.Update(windowMsg)
		m.typingModel.Update(windowMsg)
		m.resultsModel.Update(windowMsg)
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.state {
	case internal.StateMenu:
		return m.updateMenu(msg)
	case internal.StateTyping:
		return m.updateTyping(msg)
	case internal.StateResults:
		return m.updateResults(msg)
	}

	return m, nil
}

func (m *Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.menuModel.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q":
			return m, tea.Quit
		case "enter", " ":
			m.state = internal.StateTyping
			m.typingModel.Reset()
			return m, m.typingModel.Init()
		}
	}

	return m, cmd
}

func (m *Model) updateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "esc" {
			m.state = internal.StateMenu
			return m, nil
		}
	}

	_, cmd := m.typingModel.Update(msg)

	if m.typingModel.IsCompleted() {
		result := m.typingModel.GetResult()
		m.resultsModel.SetResult(result)
		m.state = internal.StateResults
		return m, nil
	}

	return m, cmd
}

func (m *Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.resultsModel.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q":
			return m, tea.Quit
		case "enter", " ", "r":
			m.state = internal.StateTyping
			m.typingModel.Reset()
			return m, m.typingModel.Init()
		case "esc":
			m.state = internal.StateMenu
			return m, nil
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	switch m.state {
	case internal.StateMenu:
		return m.menuModel.View()
	case internal.StateTyping:
		return m.typingModel.View()
	case internal.StateResults:
		return m.resultsModel.View()
	}
	return ""
}
