# DevTyper

A terminal-based typing practice tool that runs alongside long-running commands. Practice typing while waiting for your commands to complete!

![DevTyper Demo](docs/images/demo.gif)

## Why DevTyper?

Ever find yourself waiting for long commands to finish? Whether it's `npm install`, `docker build`, or creating a new project, DevTyper lets you make that time productive by practicing your typing skills.

## Features

- Practice typing while waiting for commands to complete
- Real-time WPM and accuracy tracking
- Terminal-based UI that works alongside your command output
- Smart command detection for interactive vs non-interactive commands
- Configurable word counts for practice sessions

## Quick Start

```bash
# Install
git clone https://github.com/yourusername/DevTyper.git
cd DevTyper
./install.sh

# Run with any command
devtyper docker pull nginx
devtyper npm install
```

## Documentation

- [Installation Guide](docs/installation.md)
- [User Guide](docs/user-guide.md)
- [Technical Documentation](docs/technical/README.md)
- [Contributing](docs/contributing.md)
