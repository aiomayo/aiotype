package typing

import (
	"math"
	"sync"
	"time"

	"aiotype/internal"
	"aiotype/internal/ui/shared"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	currentTest   *internal.TypingTest
	config        internal.GameConfig
	windowWidth   int
	windowHeight  int
	realTimeWPM   float64
	fadeStartTime time.Time
	mu            sync.RWMutex
}

func NewModel(config internal.GameConfig) *Model {
	test := internal.NewTest(config)
	if test == nil {
		test = internal.NewTest(internal.DefaultGameConfig())
	}
	return &Model{
		config:      config,
		currentTest: test,
	}
}

func (m *Model) Init() tea.Cmd {
	return shared.TickEvery()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil

	case shared.TickMsg:
		m.mu.Lock()
		if m.currentTest != nil && !m.currentTest.StartTime.IsZero() {
			wpm := internal.CalculateWPM(m.currentTest)
			m.realTimeWPM = wpm

			if m.isInLastFiveSeconds() {
				if m.fadeStartTime.IsZero() {
					m.fadeStartTime = time.Now()
				}
			} else {
				m.fadeStartTime = time.Time{}
			}
		}
		m.mu.Unlock()
		return m, shared.TickEvery()

	case tea.KeyMsg:
		switch msg.String() {
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
					internal.ProcessCharacter(m.currentTest, char)
				}
			}
			m.mu.Unlock()
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) IsCompleted() bool {
	return m.currentTest != nil && m.currentTest.Completed
}

func (m *Model) GetResult() *internal.TestResult {
	if m.currentTest != nil && m.currentTest.Completed {
		return internal.GenerateResult(m.currentTest)
	}
	return nil
}

func (m *Model) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	test := internal.NewTest(m.config)
	if test == nil {
		test = internal.NewTest(internal.DefaultGameConfig())
	}
	m.currentTest = test
	m.realTimeWPM = DefaultWPM
	m.fadeStartTime = time.Time{}
}

func (m *Model) isInLastFiveSeconds() bool {
	if m.currentTest == nil || m.currentTest.StartTime.IsZero() {
		return false
	}

	elapsed := time.Since(m.currentTest.StartTime)
	remaining := m.config.TestDuration - elapsed
	return remaining <= time.Duration(FadeWarningTime)*time.Second && remaining > 0
}

func (m *Model) getFadeFactor() float64 {
	if !m.isInLastFiveSeconds() {
		return 0.0
	}

	elapsed := time.Since(m.fadeStartTime)
	phase := elapsed.Seconds() / (time.Duration(FadeAnimationMs) * time.Millisecond).Seconds() * 2 * math.Pi
	return (math.Sin(phase) + 1.0) / 2.0
}

func (m *Model) getCurrentCursorStyle() lipgloss.Style {
	if !m.isInLastFiveSeconds() {
		return shared.CursorStyle
	}
	fadeFactor := m.getFadeFactor()
	interpolatedColor := shared.InterpolateColor("#ffffff", "#ff0000", fadeFactor)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color(interpolatedColor))
}
