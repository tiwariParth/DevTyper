# DevTyper Technical Documentation

## Architecture Overview

DevTyper uses a modular architecture with the following main components:

```
DevTyper/
├── cmd/                 # Command line interface
├── game/               # Core game engine
├── monitor/           # Command execution & monitoring
└── docs/              # Documentation
```

### Core Components

1. **Command Monitor**
   - Handles command execution
   - Manages process lifecycle
   - Captures command output
   - Detects interactive commands

2. **Game Engine**
   - Manages typing practice
   - Tracks statistics (WPM, accuracy)
   - Handles terminal UI
   - Word generation and validation

3. **Terminal UI**
   - Real-time character feedback
   - Command output display
   - Statistics visualization
   - Interactive mode selection

## Implementation Details

### Command Monitoring

The monitor package uses PTY/TTY for command execution, allowing:
- Real-time output capture
- Interactive command detection
- Clean process termination

### Game Engine

The game package implements:
- Word generation from common English words
- Real-time typing validation
- WPM and accuracy calculation
- Terminal UI with tcell

### State Management

Game states are managed through an FSM:
```
StateWordCountSelect → StatePlaying → StateResults
         ↑                   ↓
         └───────────────────┘
```

## Configuration

Default settings can be modified in:
- Word counts: `game.wordCountOptions`
- Interactive commands: `monitor.interactiveCommands`
