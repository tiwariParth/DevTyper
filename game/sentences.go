package game

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/parth/DevTyper/languages/golang"
	"github.com/parth/DevTyper/languages/javascript"
	"github.com/parth/DevTyper/languages/rust"
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

type Language struct {
	Templates  []string
	Variables  []string
	Types      []string
	Operations []string
}

var languages = map[string]Language{
	"go": {
		Templates:  golang.Templates,
		Variables:  golang.Variables,
		Types:      golang.Types,
		Operations: golang.Operations,
	},
	"javascript": {
		Templates:  javascript.Templates,
		Variables:  javascript.Variables,
		Types:      javascript.Types,
		Operations: javascript.Operations,
	},
	"rust": {
		Templates:  rust.Templates,
		Variables:  rust.Variables,
		Types:      rust.Types,
		Operations: rust.Operations,
	},
}

type SentenceGenerator struct {
	currentSentence string
	currentLang     string
	wordCount       int
}

func NewSentenceGenerator() *SentenceGenerator {
	return &SentenceGenerator{
		currentLang: "go", // default language
	}
}

func (sg *SentenceGenerator) SetLanguage(lang string) {
	sg.currentLang = lang
	// Reset word count for generic mode
	if lang == "generic" {
		if sg.wordCount == 0 {
			sg.wordCount = 10 // default word count
		}
	}
}

func (sg *SentenceGenerator) Generate() string {
	if sg.currentLang == "generic" {
		if sg.wordCount <= 0 {
			sg.wordCount = 10
		}
		words := make([]string, sg.wordCount)
		for i := 0; i < sg.wordCount; i++ {
			words[i] = genericWords[rand.Intn(len(genericWords))]
		}
		sg.currentSentence = strings.Join(words, " ")
		return sg.currentSentence
	}

	// Only use programming templates for non-generic modes
	if lang, ok := languages[sg.currentLang]; ok {
		template := lang.Templates[rand.Intn(len(lang.Templates))]
		parts := make([]interface{}, 0)

		for i := 0; i < strings.Count(template, "%s"); i++ {
			switch rand.Intn(3) {
			case 0:
				parts = append(parts, lang.Variables[rand.Intn(len(lang.Variables))])
			case 1:
				parts = append(parts, lang.Types[rand.Intn(len(lang.Types))])
			case 2:
				parts = append(parts, lang.Operations[rand.Intn(len(lang.Operations))])
			}
		}

		sg.currentSentence = fmt.Sprintf(template, parts...)
	}
	return sg.currentSentence
}

func (sg *SentenceGenerator) SetWordCount(count int) {
	sg.wordCount = count
}
