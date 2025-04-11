package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	templates = []string{
		"func %s() %s {",
		"var %s %s = %s",
		"type %s struct { %s %s }",
		"if %s != nil { return %s }",
		"for %s := range %s {",
		"switch %s := %s.(type) {",
		"map[%s]%s{%s: %s}",
		"func (%s *%s) %s() %s {",
	}

	variables  = []string{"err", "val", "data", "result", "item", "obj", "ctx"}
	types      = []string{"string", "int", "bool", "error", "interface{}"}
	operations = []string{"nil", "true", "false", "0", "1", "\"\""}
)

type SentenceGenerator struct {
	currentSentence string
}

func NewSentenceGenerator() *SentenceGenerator {
	return &SentenceGenerator{}
}

func (sg *SentenceGenerator) Generate() string {
	template := templates[rand.Intn(len(templates))]
	parts := make([]interface{}, 0)

	for i := 0; i < strings.Count(template, "%s"); i++ {
		switch rand.Intn(3) {
		case 0:
			parts = append(parts, variables[rand.Intn(len(variables))])
		case 1:
			parts = append(parts, types[rand.Intn(len(types))])
		case 2:
			parts = append(parts, operations[rand.Intn(len(operations))])
		}
	}

	sg.currentSentence = fmt.Sprintf(template, parts...)
	return sg.currentSentence
}
