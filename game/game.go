package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/parth/DevTyper/monitor" 
)

type GameState int

const (
	StateLanguageSelect GameState = iota
	StateTimeSelect
	StatePlaying
	StateResults
	StateTaskComplete
)

type CharacterState struct {
	char    rune
	correct bool
	typed   bool
}

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
	elapsedMinutes := time.Since(s.startTime).Minutes()
	if elapsedMinutes == 0 {
		return 0
	}
	return float64(s.wordsTyped) / elapsedMinutes
}

func (s *Stats) calculateAccuracy() float64 {
	if s.totalStrokes == 0 {
		return 100
	}
	return float64(s.totalStrokes-s.errorStrokes) / float64(s.totalStrokes) * 100
}

type Results struct {
	Language    string
	Duration    int
	WPM         float64
	Accuracy    float64
	WordsTyped  int
	TotalErrors int
}

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
	currentChars     []CharacterState
	timerDone        chan bool
	results          *Results
	taskDone         chan bool
	ForceExit        bool
	taskDescription  string
	task             *monitor.Task
}

func New(taskDone chan bool, description string, task *monitor.Task) (*Game, error) {
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
		currentChars:     make([]CharacterState, 0),
		timerDone:        make(chan bool),
		results:          &Results{},
		taskDone:         taskDone,
		ForceExit:        false,
		taskDescription:  description,
		task:             task,
	}
	game.updateCurrentChars()

	return game, nil
}

func (g *Game) updateCurrentChars() {
	g.currentChars = make([]CharacterState, len(g.currentSentence))
	for i, c := range g.currentSentence {
		g.currentChars[i] = CharacterState{char: c, correct: false, typed: false}
	}
}

func (g *Game) startTimer() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		endTime := time.Now().Add(time.Duration(g.timeOptions[g.selectedTime]) * time.Second)

		for {
			select {
			case <-ticker.C:
				remaining := time.Until(endTime).Seconds()
				g.timeRemaining = int(remaining)

				if remaining <= 0 {
					g.timerActive = false
					g.saveResults()
					g.timerDone <- true
					return
				}
			}
		}
	}()
}

func (g *Game) saveResults() {
	g.results = &Results{
		Language:    g.languages[g.selectedLanguage],
		Duration:    g.timeOptions[g.selectedTime],
		WPM:         g.stats.calculateWPM(),
		Accuracy:    g.stats.calculateAccuracy(),
		WordsTyped:  g.stats.wordsTyped,
		TotalErrors: g.stats.errorStrokes,
	}
}

func (g *Game) Run() {
gameLoop:
	for g.isRunning {
		select {
		case <-g.timerDone:
			break gameLoop
		case <-g.taskDone:
			if g.ForceExit {
				break gameLoop
			}
			g.showTaskComplete()
			// Add timeout for task completion
			go func() {
				time.Sleep(time.Second * 3)
				g.isRunning = false
			}()
		default:
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

	// Ensure clean exit
	g.screen.Fini()
	if g.results != nil {
		PrintResults(*g.results)
	}
}

func (g *Game) showTaskComplete() {
	g.screen.Beep()
	g.state = StateTaskComplete
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
			g.startTimer()
			g.stats = NewStats() // Reset stats
			g.state = StatePlaying
		}
	}
}

func (g *Game) handleInput() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.saveResults()
			g.isRunning = false
			return // Immediate return on ESC
		case tcell.KeyRune:
			if len(g.userInput) < len(g.currentSentence) {
				g.userInput += string(ev.Rune())
				g.stats.totalStrokes++
				pos := len(g.userInput) - 1
				g.currentChars[pos].typed = true
				g.currentChars[pos].correct = g.userInput[pos] == g.currentSentence[pos]
				if !g.currentChars[pos].correct {
					g.stats.errorStrokes++
				}
			}
		case tcell.KeyEnter:
			g.checkWord()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(g.userInput) > 0 {
				pos := len(g.userInput) - 1
				g.currentChars[pos].typed = false
				g.currentChars[pos].correct = false
				g.userInput = g.userInput[:pos]
			}
		}
	}
}

func (g *Game) checkWord() {
	if g.userInput == g.currentSentence {
		g.stats.wordsTyped += len(strings.Fields(g.currentSentence))
		g.currentSentence = g.sentenceGen.Generate()
		g.userInput = ""
		g.updateCurrentChars()
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

		// Draw current sentence with character-by-character coloring
		x := 7 // Starting after "Type: "
		drawText(g.screen, 1, 3, style, "Type: ")
		for i, cs := range g.currentChars {
			charStyle := style
			if cs.typed {
				if cs.correct {
					charStyle = charStyle.Foreground(tcell.ColorGreen)
				} else {
					charStyle = charStyle.Foreground(tcell.ColorRed)
				}
			}
			drawText(g.screen, x+i, 3, charStyle, string(cs.char))
		}

		// Draw user input
		drawText(g.screen, 1, 4, style, "Your input: "+g.userInput)

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

	case StateTaskComplete:
		drawText(g.screen, 1, 1, style.Bold(true), "Task Completed!")
		drawText(g.screen, 1, 3, style, g.taskDescription+" has finished")
		drawText(g.screen, 1, 5, style, "Press ESC to exit or ENTER to continue typing")
	}

	// Add task status line at bottom
	if g.taskDescription != "" {
		statusStyle := style.Bold(true)
		_, height := g.screen.Size()
		drawText(g.screen, 1, height-1, statusStyle, "Background: "+g.taskDescription)
	}

	// Show task output in bottom half of screen
	if g.task != nil {
		_, height := g.screen.Size()
		outputY := height - 10 // Reserve bottom 10 lines for output
		lines := strings.Split(g.task.GetOutput(), "\n")
		for i, line := range lines[max(0, len(lines)-8):] { // Show last 8 lines
			drawText(g.screen, 1, outputY+i, style.Foreground(tcell.ColorYellow), line)
		}
	}

	g.screen.Show()
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func PrintResults(results Results) {
	border := strings.Repeat("=", 50)
	fmt.Println(border)
	fmt.Println("🎯 DevTyper Results")
	fmt.Println(border)
	fmt.Printf("Language: %s\n", results.Language)
	fmt.Printf("Duration: %d seconds\n", results.Duration)
	fmt.Printf("Words per minute: %.1f\n", results.WPM)
	fmt.Printf("Accuracy: %.1f%%\n", results.Accuracy)
	fmt.Printf("Total words typed: %d\n", results.WordsTyped)
	fmt.Printf("Total errors: %d\n", results.TotalErrors)
	fmt.Println(border)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
