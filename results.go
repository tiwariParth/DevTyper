package main

import (
	"fmt"
	"strings"
)

type GameResults struct {
    Language    string
    Duration    int
    WPM         float64
    Accuracy    float64
    WordsTyped  int
    TotalErrors int
}

func PrintResults(results GameResults) {
    border := strings.Repeat("=", 50)
    fmt.Println(border)
    fmt.Println("🎯 DevTyper Results")
    fmt.Println(border)
    fmt.Printf("Language: %s\n", strings.ToUpper(results.Language))
    fmt.Printf("Duration: %d seconds\n", results.Duration)
    fmt.Printf("Words Per Minute: %.1f\n", results.WPM)
    fmt.Printf("Accuracy: %.1f%%\n", results.Accuracy)
    fmt.Printf("Total Words Typed: %d\n", results.WordsTyped)
    fmt.Printf("Total Errors: %d\n", results.TotalErrors)
    fmt.Println(border)
}
