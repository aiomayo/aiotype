package internal

import "time"

func CountCorrectChars(test *TypingTest) int {
	if test == nil {
		return 0
	}
	correctChars := 0
	for _, typedChar := range test.TypedChars {
		if typedChar.IsCorrect {
			correctChars++
		}
	}
	return correctChars
}

func CalculateWPM(test *TypingTest) float64 {
	if test == nil || test.StartTime.IsZero() {
		return 0
	}

	var duration time.Duration
	if test.Completed {
		duration = test.Duration
	} else {
		duration = time.Since(test.StartTime)
	}

	if duration.Seconds() <= 0 {
		return 0
	}

	correctChars := CountCorrectChars(test)

	minutes := duration.Minutes()
	return (float64(correctChars) / 5.0) / minutes
}

func CalculateAccuracy(test *TypingTest) float64 {
	if test == nil {
		return 0
	}
	totalChars := len(test.TypedChars)
	if totalChars == 0 {
		return 100.0
	}

	correctChars := CountCorrectChars(test)

	return (float64(correctChars) / float64(totalChars)) * 100.0
}

func GenerateResult(test *TypingTest) *TestResult {
	if test == nil {
		return nil
	}
	totalChars := len(test.TypedChars)
	correctChars := CountCorrectChars(test)

	correctWords := 0
	if len(test.TypedChars) > 0 {
		wordStart := 0
		for i, typedChar := range test.TypedChars {
			if typedChar.Character == ' ' {
				wordIsCorrect := true
				for j := wordStart; j <= i; j++ {
					if !test.TypedChars[j].IsCorrect {
						wordIsCorrect = false
						break
					}
				}
				if wordIsCorrect {
					correctWords++
				}
				wordStart = i + 1
			}
		}
	}

	return &TestResult{
		WPM:          CalculateWPM(test),
		Accuracy:     CalculateAccuracy(test),
		TotalWords:   len(test.Words),
		CorrectWords: correctWords,
		TotalChars:   totalChars,
		CorrectChars: correctChars,
		ErrorCount:   totalChars - correctChars,
		TestDuration: test.Duration,
		CompletedAt:  test.EndTime,
	}
}
