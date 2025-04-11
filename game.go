package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
	"time"
)

type Game struct {
	screen           tcell.Screen
	sentenceGen      *SentenceGenerator
	currentSentence  string
	userInput        string
	isRunning        bool
	stats            *Stats
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
		screen:          screen,
		sentenceGen:     sentenceGen,
		currentSentence: sentenceGen.Generate(),
		isRunning:       true,
		stats:           NewStats(),
	}

	return game, nil
}

func (g *Game) Run() {
	g.stats.startTime = time.Now()

	for g.isRunning {
		g.handleInput()
		g.draw()
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

func (g *Game) draw() {
	g.screen.Clear()
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	// Draw title
	drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Developer Typing Practice")

	// Draw current sentence
	drawText(g.screen, 1, 3, style, "Type: "+g.currentSentence)

	// Update the input validation logic
	inputStyle := style
	if len(g.userInput) > 0 {
		targetLen := len(g.userInput)
		if targetLen <= len(g.currentSentence) {
			// Compare only the typed portion
			expected := g.currentSentence[:targetLen]
			if g.userInput == expected {
				inputStyle = inputStyle.Foreground(tcell.ColorGreen)
			} else {
				inputStyle = inputStyle.Foreground(tcell.ColorRed)
			}
		} else {
			// User has typed more characters than the target sentence
			inputStyle = inputStyle.Foreground(tcell.ColorRed)
		}
	}
	drawText(g.screen, 1, 4, inputStyle, "Your input: "+g.userInput)

	// Draw stats
	drawText(g.screen, 1, 6, style, fmt.Sprintf("WPM: %.1f", g.stats.calculateWPM()))
	drawText(g.screen, 1, 7, style, fmt.Sprintf("Accuracy: %.1f%%", g.stats.calculateAccuracy()))
	drawText(g.screen, 1, 8, style, fmt.Sprintf("Words typed: %d", g.stats.wordsTyped))

	drawText(g.screen, 1, 10, style, "Press ESC to exit")

	g.screen.Show()
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
