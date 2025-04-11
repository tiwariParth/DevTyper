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
	StateWordCountSelect
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
	cursorX          int
	cursorY          int
	lastOutput       []string
	outputStartRow   int
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
		cursorX:           7,
		cursorY:           3,
		lastOutput:       []string{},
		outputStartRow:   0,
	}
	return game, nil
}

func (g *Game) updateCurrentChars() {
	g.currentChars = make([]CharacterState, len(g.currentSentence))
	for i, c := range g.currentSentence {
		g.currentChars[i] = CharacterState{char: c, correct: false, typed: false}
	}
}

func (g *Game) updateCommandOutput() {
	if g.task == nil {
		return
	}

	// Get up to 5 recent lines from the task's buffer
	g.lastOutput = g.task.GetRecentOutput(5)
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
			case StateWordCountSelect:
				g.handleWordCountSelect()
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
				g.state = StateWordCountSelect
			} else {
				g.isRunning = false // Just wait for task
			}
		}
	}
}

func (g *Game) handleWordCountSelect() {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.state = StateMode
		case tcell.KeyUp:
			g.selectedCount = (g.selectedCount - 1 + len(g.wordCountOptions)) % len(g.wordCountOptions)
		case tcell.KeyDown:
			g.selectedCount = (g.selectedCount + 1) % len(g.wordCountOptions)
		case tcell.KeyEnter:
			g.sentenceGen.SetWordCount(g.wordCountOptions[g.selectedCount])
			g.currentSentence = g.sentenceGen.Generate()
			g.updateCurrentChars()
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
	width, height := g.screen.Size()

	// Update command output before drawing
	g.updateCommandOutput()

	// Calculate layout more carefully
	commandOutputHeight := min(len(g.lastOutput)+2, 7) // +2 for border and title, max 7 lines total

	// Clear screen before drawing
	g.screen.Fill(' ', style)

	// Calculate layout
	availableHeight := height - commandOutputHeight - 1 // -1 for bottom status line

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

	case StateWordCountSelect:
		drawText(g.screen, 1, 1, style.Bold(true), "DevTyper - Select Word Count")
		drawText(g.screen, 1, 3, style, "Use Up/Down arrows to select, Enter to confirm:")

		for i, count := range g.wordCountOptions {
			countStyle := style
			if i == g.selectedCount {
				countStyle = countStyle.Background(tcell.ColorBlue)
			}
			drawText(g.screen, 3, 5+i, countStyle, fmt.Sprintf("%d words", count))
		}

	case StatePlaying:
		// Draw header with border
		drawBorder(g.screen, 0, 0, width-1, 2, style)
		drawText(g.screen, 2, 1, style.Bold(true), "DevTyper - Typing Practice")

		// Calculate how much space we have for text area
		textAreaHeight := availableHeight - 8 // Reserve space for stats box

		// Draw text area with border
		drawBorder(g.screen, 0, 2, width-1, textAreaHeight+2, style)

		// Draw sentence with wrapping and cursor
		maxWidth := width - 8
		wrappedSentence := wrapText(g.currentSentence, maxWidth)
		currentPos := len(g.userInput)

		drawText(g.screen, 2, 3, style, "Type:")

		for i, line := range wrappedSentence {
			x := 7
			lineStart := 0
			if i > 0 {
				for _, prevLine := range wrappedSentence[:i] {
					lineStart += len(prevLine)
				}
			}

			// Draw each character in the line
			for j, char := range line {
				pos := lineStart + j
				charStyle := style

				if pos == currentPos {
					g.cursorX = x + j
					g.cursorY = 3 + i
					charStyle = charStyle.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
				} else if pos < len(g.currentChars) {
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

			// Show cursor at end of input if needed
			if currentPos == len(g.currentSentence) && i == len(wrappedSentence)-1 {
				g.cursorX = x + len(line)
				g.cursorY = 3 + i
			}
		}

		// Draw stats box
		statsY := textAreaHeight + 4
		drawBorder(g.screen, 0, statsY-1, width-1, statsY+3, style)
		drawText(g.screen, 2, statsY, style, fmt.Sprintf("WPM: %.1f | Accuracy: %.1f%% | Words: %d",
			g.stats.calculateWPM(),
			g.stats.calculateAccuracy(),
			g.stats.wordsTyped))
		drawText(g.screen, 2, statsY+2, style, "Press ESC to exit")

		// Set the starting row for command output
		g.outputStartRow = statsY + 5

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

	// Draw command output in dedicated area at the bottom with border
	outputY := height - commandOutputHeight - 1 // -1 for status line

	if len(g.lastOutput) > 0 && g.state != StateTaskComplete {
		g.outputStartRow = outputY

		// Draw output box
		outputStyle := style.Foreground(tcell.ColorYellow).Background(tcell.ColorReset)
		drawBorder(g.screen, 0, outputY, width-1, height-2, style)
		drawText(g.screen, 2, outputY, style.Bold(true), "Command Output")

		// Draw command output with better formatting
		for i, line := range g.lastOutput {
			if i >= 5 { // Show max 5 lines
				break
			}
			// Trim line if needed
			if len(line) > width-6 {
				line = line[:width-9] + "..."
			}
			drawText(g.screen, 2, outputY+i+1, outputStyle, line)
		}
	}

	// Draw a clear status line at the very bottom with border
	statusY := height - 1
	drawText(g.screen, 1, statusY, style.Bold(true),
		fmt.Sprintf("Task: %s | Press ESC to exit", g.taskDescription))

	// Show cursor
	g.screen.ShowCursor(g.cursorX, g.cursorY)
	g.screen.Show()
}

// Add helper function for drawing borders
func drawBorder(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	// Draw corners
	s.SetContent(x1, y1, '┌', nil, style)
	s.SetContent(x2, y1, '┐', nil, style)
	s.SetContent(x1, y2, '└', nil, style)
	s.SetContent(x2, y2, '┘', nil, style)

	// Draw horizontal lines
	for x := x1 + 1; x < x2; x++ {
		s.SetContent(x, y1, '─', nil, style)
		s.SetContent(x, y2, '─', nil, style)
	}

	// Draw vertical lines
	for y := y1 + 1; y < y2; y++ {
		s.SetContent(x1, y, '│', nil, style)
		s.SetContent(x2, y, '│', nil, style)
	}
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

// Helper for calculating minimum
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
