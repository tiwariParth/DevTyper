#!/bin/bash

echo "Installing DevTyper..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

# Build the binary
go build -o devtyper cmd/devtyper/main.go

# Create installation directory
sudo mkdir -p /usr/local/bin

# Move binary to installation directory
sudo mv devtyper /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/devtyper

echo "DevTyper installed successfully! Try running: devtyper"
