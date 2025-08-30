package internal

import "time"

type GameState int

const (
	StateMenu GameState = iota
	StateTyping
	StateResults
)

type TypedChar struct {
	Character rune
	IsCorrect bool
	Timestamp time.Time
}

type WordStatus struct {
	StartIndex int
	EndIndex   int
	HasError   bool
	IsComplete bool
}

type TypingTest struct {
	Words        []string
	TargetText   string
	TypedChars   []TypedChar
	CurrentPos   int
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Completed    bool
	WordStatuses []WordStatus
}

type TestResult struct {
	WPM          float64
	Accuracy     float64
	TotalWords   int
	CorrectWords int
	TotalChars   int
	CorrectChars int
	ErrorCount   int
	TestDuration time.Duration
	CompletedAt  time.Time
}

type GameConfig struct {
	TestDuration time.Duration
	WordCount    int
}
