package internal

import (
	"math/rand"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var commonWords = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "i",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "these", "give", "day", "most", "us",
	"is", "water", "long", "very", "still", "through", "down",
	"may", "such", "here", "were", "been", "much",
	"where", "too", "each", "many",
	"has", "more", "life", "should",
	"being", "made", "before", "might", "did", "every", "large", "often",
	"together", "asked", "house", "don't", "world", "going", "school", "important", "until", "form",
	"food", "keep", "children", "feet", "land", "side", "without", "boy", "once", "animal",
	"enough", "took", "sometimes", "four", "head", "above", "kind", "began", "almost",
	"live", "page", "got", "earth", "need", "far", "hand", "high", "mother",
	"light", "country", "father", "let", "night", "picture", "study", "second", "soon",
	"story", "since", "white", "ever", "paper", "hard", "near", "sentence", "better", "best",
	"across", "during", "today", "however", "sure", "knew", "it's", "try", "told", "young",
	"sun", "thing", "whole", "hear", "example", "heard", "several", "change", "answer", "room",
	"sea", "against", "top", "turned", "learn", "point", "city", "play", "toward", "five",
	"himself", "usually", "money", "seen", "didn't", "car", "morning", "i'm", "body", "upon",
	"family", "later", "turn", "move", "face", "door", "cut", "done", "group", "true",
	"leave", "color", "red", "friend", "pretty", "eat", "front", "feel", "fact",
	"week", "eye", "same", "another", "left", "call", "while", "right",
	"find", "part", "place", "under", "name", "help", "low", "line", "cause", "mean",
	"differ", "old", "tell", "follow", "around", "three", "small", "set", "put", "end",
	"why", "again", "off", "went", "number", "men", "found", "between", "home",
	"big", "air", "own", "read", "last", "along", "next", "below", "saw",
	"something", "thought", "both", "few", "those", "always", "looked",
}

func NewTest(config GameConfig) *TypingTest {
	words := make([]string, config.WordCount)
	for i := 0; i < config.WordCount; i++ {
		words[i] = commonWords[rng.Intn(len(commonWords))]
	}

	targetText := strings.Join(words, " ")

	wordStatuses := []WordStatus{}
	currentIndex := 0
	for i, word := range words {
		wordStatuses = append(wordStatuses, WordStatus{
			StartIndex: currentIndex,
			EndIndex:   currentIndex + len(word) - 1,
			HasError:   false,
			IsComplete: false,
		})
		currentIndex += len(word)

		if i < len(words)-1 {
			wordStatuses = append(wordStatuses, WordStatus{
				StartIndex: currentIndex,
				EndIndex:   currentIndex,
				HasError:   false,
				IsComplete: false,
			})
			currentIndex++
		}
	}

	return &TypingTest{
		Words:        words,
		TargetText:   targetText,
		TypedChars:   make([]TypedChar, 0, len(targetText)),
		CurrentPos:   0,
		Completed:    false,
		WordStatuses: wordStatuses,
	}
}

func ProcessCharacter(test *TypingTest, char rune) bool {
	if test == nil || test.Completed || test.CurrentPos >= len(test.TargetText) {
		return false
	}

	if test.StartTime.IsZero() {
		test.StartTime = time.Now()
	}

	if test.CurrentPos < 0 || test.CurrentPos >= len(test.TargetText) {
		return false
	}

	expected := rune(test.TargetText[test.CurrentPos])
	isCorrect := char == expected

	test.TypedChars = append(test.TypedChars, TypedChar{
		Character: char,
		IsCorrect: isCorrect,
		Timestamp: time.Now(),
	})
	test.CurrentPos++

	updateWordStatus(test)

	if test.CurrentPos >= len(test.TargetText) {
		test.Completed = true
		test.EndTime = time.Now()
		test.Duration = test.EndTime.Sub(test.StartTime)
		return true
	}

	return false
}

func ProcessBackspace(test *TypingTest) {
	if test == nil || len(test.TypedChars) == 0 {
		return
	}

	prevPos := test.CurrentPos

	test.TypedChars = test.TypedChars[:len(test.TypedChars)-1]
	test.CurrentPos = len(test.TypedChars)

	if test.CurrentPos < 0 {
		test.CurrentPos = 0
	}

	updateWordStatusOnBackspace(test, prevPos)
}

func GetWordIndexForPosition(test *TypingTest, position int) int {
	for i, ws := range test.WordStatuses {
		if position >= ws.StartIndex && position <= ws.EndIndex {
			return i
		}
	}
	return -1
}

func updateWordStatus(test *TypingTest) {
	if test == nil || test.CurrentPos == 0 {
		return
	}

	charPos := test.CurrentPos - 1

	unitIndex := GetWordIndexForPosition(test, charPos)
	if unitIndex == -1 {
		return
	}

	if charPos == test.WordStatuses[unitIndex].EndIndex {
		test.WordStatuses[unitIndex].IsComplete = true

		hasError := false
		for i := test.WordStatuses[unitIndex].StartIndex; i <= test.WordStatuses[unitIndex].EndIndex; i++ {
			if i < len(test.TypedChars) && !test.TypedChars[i].IsCorrect {
				hasError = true
				break
			}
		}
		test.WordStatuses[unitIndex].HasError = hasError
	}
}

func updateWordStatusOnBackspace(test *TypingTest, previousPos int) {
	if test == nil {
		return
	}

	unitIndex := GetWordIndexForPosition(test, previousPos-1)
	if unitIndex == -1 {
		return
	}

	if test.CurrentPos <= test.WordStatuses[unitIndex].EndIndex {
		test.WordStatuses[unitIndex].IsComplete = false
		test.WordStatuses[unitIndex].HasError = false
	}
}
