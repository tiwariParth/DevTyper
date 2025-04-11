package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/parth/DevTyper/languages/golang"
	"github.com/parth/DevTyper/languages/javascript"
	"github.com/parth/DevTyper/languages/rust"
)

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
}

func NewSentenceGenerator() *SentenceGenerator {
	return &SentenceGenerator{
		currentLang: "go", // default language
	}
}

func (sg *SentenceGenerator) SetLanguage(lang string) {
	if _, ok := languages[lang]; ok {
		sg.currentLang = lang
	}
}

func (sg *SentenceGenerator) Generate() string {
	lang := languages[sg.currentLang]
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
	return sg.currentSentence
}
