package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameState int

const (
	StateLanguageSelect GameState = iota
	StateTimeSelect
	StatePlaying
	StateResults
)

type Game struct {
	screen           tcell.Screen
	sentenceGen      *SentenceGenerator
	currentSentence  string
	userInput        string
	isRunning        bool
	stats            *Stats
	state            GameState
	selectedLanguage int
	languages        []string
	timeOptions      []int
	selectedTime     int
	timeRemaining    int
	timerActive      bool
	lastTick         time.Time
}

func NewGame() (*Game, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	sentenceGen := NewSentenceGenerator()

	game := &Game{
		screen:           screen,
		sentenceGen:      sentenceGen,
		currentSentence:  sentenceGen.Generate(),
		isRunning:        true,
		stats:            NewStats(),
		state:            StateLanguageSelect,
		selectedLanguage: 0,
		languages:        []string{"go", "javascript", "rust"},
		timeOptions:      []int{30, 60, 90, 120},
		selectedTime:     0,
		timeRemaining:    0,
		timerActive:      false,
		lastTick:         time.Now(),
	}

	return game, nil
}

func (g *Game) Run() {
	for g.isRunning {
		if g.timerActive {
			g.updateTimer()
		}
		switch g.state {
		case StateLanguageSelect:
			g.handleLanguageSelect()
		case StateTimeSelect:
			g.handleTimeSelect()
		case StatePlaying:
			g.handleInput()
		case StateResults:
			g.handleResults()
		}
		g.draw()
	}
}

func (g *Game) handleLanguageSelect() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.isRunning = false
		case tcell.KeyUp:
			g.selectedLanguage = (g.selectedLanguage - 1 + len(g.languages)) % len(g.languages)
		case tcell.KeyDown:
			g.selectedLanguage = (g.selectedLanguage + 1) % len(g.languages)
		case tcell.KeyEnter:
			g.state = StateTimeSelect
		}
	}
}

func (g *Game) handleTimeSelect() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.state = StateLanguageSelect
		case tcell.KeyUp:
			g.selectedTime = (g.selectedTime - 1 + len(g.timeOptions)) % len(g.timeOptions)
		case tcell.KeyDown:
			g.selectedTime = (g.selectedTime + 1) % len(g.timeOptions)
		case tcell.KeyEnter:
			g.timeRemaining = g.timeOptions[g.selectedTime]
			g.timerActive = true
			g.lastTick = time.Now()
			g.stats = NewStats() // Reset stats
			g.state = StatePlaying
		}
	}
}

func (g *Game) updateTimer() {
	now := time.Now()
	elapsed := now.Sub(g.lastTick).Seconds()
	g.lastTick = now

	g.timeRemaining -= int(elapsed)
	if g.timeRemaining <= 0 {
		g.timerActive = false
		g.state = StateResults
	}
}

func (g *Game) handleInput() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.isRunning = false
		case tcell.KeyRune:
			g.userInput += string(ev.Rune())
			g.stats.totalStrokes++
			if len(g.userInput) <= len(g.currentSentence) &&
				g.userInput[len(g.userInput)-1] != g.currentSentence[len(g.userInput)-1] {
				g.stats.errorStrokes++
			}
		case tcell.KeyEnter:
			g.checkWord()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(g.userInput) > 0 {
				g.userInput = g.userInput[:len(g.userInput)-1]
			}
		}
	}
}

func (g *Game) checkWord() {
	if g.userInput == g.currentSentence {
		g.stats.wordsTyped += len(strings.Fields(g.currentSentence))
		g.currentSentence = g.sentenceGen.Generate()
		g.userInput = ""
	}
}

func (g *Game) handleResults() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			g.state = StateLanguageSelect
		}
	}
}

func (g *Game) draw() {
	g.screen.Clear()
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	switch g.state {
	case StateLanguageSelect:
		drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Select Programming Language")
		drawText(g.screen, 1, 3, style, "Use Up/Down arrows to select, Enter to confirm:")

		for i, lang := range g.languages {
			langStyle := style
			if i == g.selectedLanguage {
				langStyle = langStyle.Background(tcell.ColorBlue)
			}
			drawText(g.screen, 3, 5+i, langStyle, lang)
		}

	case StateTimeSelect:
		drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Select Time Limit")
		drawText(g.screen, 1, 3, style, "Use Up/Down arrows to select, Enter to confirm:")

		for i, t := range g.timeOptions {
			timeStyle := style
			if i == g.selectedTime {
				timeStyle = timeStyle.Background(tcell.ColorBlue)
			}
			drawText(g.screen, 3, 5+i, timeStyle, fmt.Sprintf("%d seconds", t))
		}
		drawText(g.screen, 1, 6+len(g.timeOptions), style, "Press ESC to go back")

	case StatePlaying:
		drawText(g.screen, 1, 1, style.Bold(true),
			fmt.Sprintf("DevTyper - %s Practice", strings.ToUpper(g.languages[g.selectedLanguage])))

		drawText(g.screen, 1, 2, style, fmt.Sprintf("Time remaining: %d seconds", g.timeRemaining))

		drawText(g.screen, 1, 3, style, "Type: "+g.currentSentence)

		inputStyle := style
		if len(g.userInput) > 0 {
			targetLen := len(g.userInput)
			if targetLen <= len(g.currentSentence) {
				expected := g.currentSentence[:targetLen]
				if g.userInput == expected {
					inputStyle = inputStyle.Foreground(tcell.ColorGreen)
				} else {
					inputStyle = inputStyle.Foreground(tcell.ColorRed)
				}
			} else {
				inputStyle = inputStyle.Foreground(tcell.ColorRed)
			}
		}
		drawText(g.screen, 1, 4, inputStyle, "Your input: "+g.userInput)

		drawText(g.screen, 1, 6, style, fmt.Sprintf("WPM: %.1f", g.stats.calculateWPM()))
		drawText(g.screen, 1, 7, style, fmt.Sprintf("Accuracy: %.1f%%", g.stats.calculateAccuracy()))
		drawText(g.screen, 1, 8, style, fmt.Sprintf("Words typed: %d", g.stats.wordsTyped))

		drawText(g.screen, 1, 10, style, "Press ESC to exit")

	case StateResults:
		drawText(g.screen, 1, 1, style.Bold(true), "Time's Up! - Final Results")
		drawText(g.screen, 1, 3, style, fmt.Sprintf("Words per minute: %.1f", g.stats.calculateWPM()))
		drawText(g.screen, 1, 4, style, fmt.Sprintf("Accuracy: %.1f%%", g.stats.calculateAccuracy()))
		drawText(g.screen, 1, 5, style, fmt.Sprintf("Total words typed: %d", g.stats.wordsTyped))
		drawText(g.screen, 1, 7, style, "Press Enter/ESC to return to language selection")
	}

	g.screen.Show()
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
