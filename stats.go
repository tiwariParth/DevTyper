package main

import (
	"time"
)

type Stats struct {
	startTime    time.Time
	wordsTyped   int
	totalStrokes int
	errorStrokes int
}

func NewStats() *Stats {
	return &Stats{
		startTime: time.Now(),
	}
}

func (s *Stats) calculateWPM() float64 {
	elapsed := time.Since(s.startTime).Minutes()
	if elapsed == 0 {
		return 0
	}
	return float64(s.wordsTyped) / elapsed
}

func (s *Stats) calculateAccuracy() float64 {
	if s.totalStrokes == 0 {
		return 100.0
	}
	correct := s.totalStrokes - s.errorStrokes
	return (float64(correct) / float64(s.totalStrokes)) * 100
}
