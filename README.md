# DevTyper

A terminal-based typing practice game for developers, designed to run alongside long-running commands.

## Installation

### Option 1: Direct Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/DevTyper.git
cd DevTyper

# Run installation script
chmod +x install.sh
./install.sh
```

### Option 2: Manual Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/DevTyper.git
cd DevTyper

# Build the binary
go build -o devtyper cmd/devtyper/main.go

# Move to path (optional)
sudo mv devtyper /usr/local/bin/
```

## Usage

### Basic Usage
```bash
# Run with a command
devtyper docker pull nginx

# Run with force exit on task completion
devtyper -force-exit npm install

# Run with multi-word commands
devtyper "eksctl create cluster --name my-cluster"
```

### Common Examples
1. Docker Image Pull:
```bash
devtyper docker pull ubuntu:latest
```

2. NPM Install:
```bash
devtyper npm install
```

3. Kubernetes Cluster Creation:
```bash
devtyper "eksctl create cluster --name test-cluster --nodes 3"
```

4. Next.js Project Creation:
```bash
devtyper npx create-next-app my-app
```

### Game Controls
- Use arrow keys to navigate menus
- Enter to select
- ESC to exit
- Enter to submit typed text
- Backspace to correct mistakes

### Features
- Multiple programming language support (Go, JavaScript, Rust)
- Real-time typing feedback
- WPM and accuracy tracking
- Background task monitoring
- Configurable time limits

## Development

### Project Structure
```
DevTyper/
├── cmd/
│   └── devtyper/
│       └── main.go
├── game/
│   ├── game.go
│   └── sentences.go
├── languages/
│   ├── golang/
│   ├── javascript/
│   └── rust/
├── monitor/
│   ├── detector.go
│   └── task.go
└── go.mod
```

### Building from Source
```bash
go mod tidy
go build -o devtyper cmd/devtyper/main.go
```
