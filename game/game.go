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
	StateMode GameState = iota
	StatePlaying
	StateResults
	StateTaskComplete
	StateError
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
	wordCountOptions []int
	selectedCount    int
	currentChars     []CharacterState
	results          *Results
	taskDone         chan bool
	ForceExit        bool
	taskDescription  string
	task             *monitor.Task
	modeOptions      []string
	selectedMode     int
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
		isRunning:        true,
		stats:            NewStats(),
		state:            StateMode,
		wordCountOptions: []int{10, 25, 50, 100},
		selectedCount:    0,
		currentChars:     make([]CharacterState, 0),
		results:          &Results{},
		taskDone:         taskDone,
		ForceExit:        false,
		taskDescription:  description,
		task:             task,
		modeOptions:      []string{"Practice Typing", "Wait for Task"},
		selectedMode:     0,
	}
	return game, nil
}

func (g *Game) updateCurrentChars() {
	g.currentChars = make([]CharacterState, len(g.currentSentence))
	for i, c := range g.currentSentence {
		g.currentChars[i] = CharacterState{char: c, correct: false, typed: false}
	}
}

func (g *Game) Run() {
gameLoop:
	for g.isRunning {
		select {
		case <-g.taskDone:
			if g.ForceExit || g.task.IsComplete() {
				g.showTaskComplete()
				break gameLoop
			}
			g.showTaskComplete()
		default:
			switch g.state {
			case StateMode:
				g.handleModeSelect()
			case StatePlaying:
				g.handleInput()
			case StateResults:
				g.handleResults()
			case StateError:
				g.draw()
			case StateTaskComplete:
				g.draw()
			}
			g.draw()
		}
	}

	// Ensure clean exit
	g.Cleanup()
}

func (g *Game) showTaskComplete() {
	g.screen.Beep()
	g.state = StateTaskComplete
	g.screen.Sync()

	// Draw completion message
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	drawText(g.screen, 1, 1, style.Bold(true), "Task Completed!")
	if g.task.HasError() {
		drawText(g.screen, 1, 3, style.Foreground(tcell.ColorRed),
			fmt.Sprintf("Error: %s", g.task.GetError()))
	} else {
		drawText(g.screen, 1, 3, style, "Task completed successfully!")
	}
	drawText(g.screen, 1, 5, style, "Press ESC to exit")
	g.screen.Show()
}

func (g *Game) showError(errMsg string) {
	g.screen.Beep()
	g.state = StateError
	g.taskDescription = errMsg
	g.screen.Sync()
}

func (g *Game) handleModeSelect() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.isRunning = false
		case tcell.KeyUp, tcell.KeyDown:
			g.selectedMode = (g.selectedMode + 1) % len(g.modeOptions)
		case tcell.KeyEnter:
			if g.selectedMode == 0 {
				g.state = StatePlaying
				g.currentSentence = g.sentenceGen.Generate()
				g.updateCurrentChars()
			} else {
				g.isRunning = false // Just wait for task
			}
		}
	}
}

func (g *Game) handleInput() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.Cleanup()
			g.isRunning = false
			return
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

func (g *Game) saveResults() {
	g.results = &Results{
		WPM:         g.stats.calculateWPM(),
		Accuracy:    g.stats.calculateAccuracy(),
		WordsTyped:  g.stats.wordsTyped,
		TotalErrors: g.stats.errorStrokes,
	}
}

func (g *Game) handleResults() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			g.Cleanup()
			g.isRunning = false
		}
	}
}

func (g *Game) Cleanup() {
	g.saveResults()
	g.screen.Clear()
	g.screen.Sync()
	g.screen.Fini()
	// Reset terminal state
	fmt.Print("\033[?25h") // Show cursor
	fmt.Print("\033[2J\033[H") // Clear screen
}

func (g *Game) ShowError(message string) {
	g.state = StateError
	g.taskDescription = message
	g.screen.Sync()
}

// Add new helper function to wrap text
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)
	return lines
}

func (g *Game) draw() {
	g.screen.Clear()
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	switch g.state {
	case StateMode:
		drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Choose Mode")
		drawText(g.screen, 1, 3, style, "Use Up/Down arrows to select, Enter to confirm:")

		for i, mode := range g.modeOptions {
			modeStyle := style
			if i == g.selectedMode {
				modeStyle = modeStyle.Background(tcell.ColorBlue)
			}
			drawText(g.screen, 3, 5+i, modeStyle, mode)
		}

	case StatePlaying:
		drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Typing Practice")

		// Get screen width for text wrapping
		width, _ := g.screen.Size()
		maxWidth := width - 8 // Account for left margin and "Type: " prefix
		
		// Wrap the sentence and user input
		wrappedSentence := wrapText(g.currentSentence, maxWidth)
		wrappedInput := wrapText(g.userInput, maxWidth)

		// Draw sentence with wrapping
		drawText(g.screen, 1, 3, style, "Type:")
		for i, line := range wrappedSentence {
			x := 7
			lineStart := 0
			if i > 0 {
				// Calculate how many characters we've already displayed
				for _, prevLine := range wrappedSentence[:i] {
					lineStart += len(prevLine)
				}
			}
			
			// Draw each character in the line
			for j, char := range line {
				pos := lineStart + j
				charStyle := style
				if pos < len(g.currentChars) {
					cs := g.currentChars[pos]
					if cs.typed {
						if cs.correct {
							charStyle = charStyle.Foreground(tcell.ColorGreen)
						} else {
							charStyle = charStyle.Foreground(tcell.ColorRed)
						}
					}
				}
				drawText(g.screen, x+j, 3+i, charStyle, string(char))
			}
		}

		// Draw user input with wrapping
		inputY := 4 + len(wrappedSentence)
		drawText(g.screen, 1, inputY, style, "Your input:")
		for i, line := range wrappedInput {
			drawText(g.screen, 7, inputY+i, style, line)
		}

		// Adjust stats position based on wrapped text
		statsY := inputY + len(wrappedInput) + 2
		drawText(g.screen, 1, statsY, style, fmt.Sprintf("WPM: %.1f", g.stats.calculateWPM()))
		drawText(g.screen, 1, statsY+1, style, fmt.Sprintf("Accuracy: %.1f%%", g.stats.calculateAccuracy()))
		drawText(g.screen, 1, statsY+2, style, fmt.Sprintf("Words typed: %d", g.stats.wordsTyped))
		drawText(g.screen, 1, statsY+4, style, "Press ESC to exit")

	case StateResults:
		drawText(g.screen, 1, 1, style.Bold(true), "Time's Up! - Final Results")
		drawText(g.screen, 1, 3, style, fmt.Sprintf("Words per minute: %.1f", g.stats.calculateWPM()))
		drawText(g.screen, 1, 4, style, fmt.Sprintf("Accuracy: %.1f%%", g.stats.calculateAccuracy()))
		drawText(g.screen, 1, 5, style, fmt.Sprintf("Total words typed: %d", g.stats.wordsTyped))
		drawText(g.screen, 1, 7, style, "Press Enter/ESC to exit")

	case StateTaskComplete:
		width, _ := g.screen.Size()
		task := g.task.GetOutput() // Get final output
		lines := wrapText(task, width-2) // Wrap task output

		drawText(g.screen, 1, 1, style.Bold(true), "Task Completed!")
		drawText(g.screen, 1, 3, style, g.taskDescription+" has finished")
		
		// Display final task output
		outputY := 5
		maxLines := 8 // Show last 8 lines
		start := len(lines)
		if start > maxLines {
			start = len(lines) - maxLines
		}
		for i, line := range lines[start:] {
			drawText(g.screen, 1, outputY+i, style.Foreground(tcell.ColorYellow), line)
		}
		
		drawText(g.screen, 1, outputY+maxLines+2, style, "Press ESC to exit")

		// Handle ESC immediately
		ev := g.screen.PollEvent()
		if ev, ok := ev.(*tcell.EventKey); ok {
			if ev.Key() == tcell.KeyEscape {
				g.Cleanup()
				g.isRunning = false
				return
			}
		}

	case StateError:
		drawText(g.screen, 1, 1, style.Bold(true).Foreground(tcell.ColorRed), "Error!")
		drawText(g.screen, 1, 3, style, g.taskDescription)
		drawText(g.screen, 1, 5, style, "Press ESC to exit")
		ev := g.screen.PollEvent()
		if ev, ok := ev.(*tcell.EventKey); ok {
			if ev.Key() == tcell.KeyEscape {
				g.Cleanup()
				g.isRunning = false
				return
			}
		}
	}

	if g.taskDescription != "" {
		statusStyle := style.Bold(true)
		_, height := g.screen.Size()
		drawText(g.screen, 1, height-1, statusStyle, "Background: "+g.taskDescription)
	}

	if g.task != nil {
		_, height := g.screen.Size()
		outputY := height - 10
		lines := strings.Split(g.task.GetOutput(), "\n")
		for i, line := range lines[max(0, len(lines)-8):] {
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
