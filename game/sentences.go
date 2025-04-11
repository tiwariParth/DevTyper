package game

import (
	"math/rand"
	"strings"
)

var genericWords = []string{
	// Common words
	"the", "be", "to", "of", "and", "that", "have", "with", "this", "from",
	"they", "say", "her", "she", "will", "one", "all", "would", "there", "their",
	"what", "out", "about", "who", "get", "which", "when", "make", "can", "like",
	"time", "just", "know", "take", "people", "into", "year", "good", "some", "could",
	"them", "see", "other", "than", "then", "now", "look", "only", "come", "its",
	"over", "think", "also", "back", "after", "use", "two", "how", "our", "work",
	"first", "well", "way", "even", "new", "want", "because", "any", "these", "give",
	"day", "most", "us", "should", "need", "much", "right", "without", "through", "own",
	"too", "here", "still", "such", "last", "great", "long", "small", "might", "around",
	"while", "those", "always", "world", "both", "life", "where", "next", "being", "keep",
}

type SentenceGenerator struct {
	currentSentence string
	wordCount       int
}

func NewSentenceGenerator() *SentenceGenerator {
	return &SentenceGenerator{
		wordCount: 25, // Default word count
	}
}

func (sg *SentenceGenerator) Generate() string {
	words := make([]string, sg.wordCount)
	for i := 0; i < sg.wordCount; i++ {
		words[i] = genericWords[rand.Intn(len(genericWords))]
	}
	sg.currentSentence = strings.Join(words, " ")
	return sg.currentSentence
}

func (sg *SentenceGenerator) SetWordCount(count int) {
	sg.wordCount = count
}
