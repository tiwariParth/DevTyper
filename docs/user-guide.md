# DevTyper User Guide

## Basic Usage

```bash
devtyper [options] <command>
```

### Options

- `--force-exit`: Exit game when command completes
- `--keep-alive`: Keep command running after exiting game (default: true)

### Examples

1. Basic usage:
```bash
devtyper npm install
```

2. With Docker:
```bash
devtyper docker build -t myapp .
```

3. Long-running commands:
```bash
devtyper "kubectl apply -f manifests/"
```

## Interactive Commands

Some commands require user input (like `npx create-next-app`). For these, DevTyper will suggest non-interactive alternatives:

```bash
# Instead of:
devtyper npx create-next-app

# Use:
devtyper npx create-next-app my-app --typescript --tailwind
```

## Game Controls

1. Mode Selection:
   - Up/Down: Select mode
   - Enter: Confirm
   - ESC: Exit

2. Typing Practice:
   - Type the shown text
   - Backspace to correct
   - ESC to exit
   - Real-time feedback with colors
