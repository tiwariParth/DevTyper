#!/bin/bash

echo "Installing DevTyper..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

# Clean any existing binary
sudo rm -f /usr/local/bin/devtyper

# Get dependencies
go mod tidy

# Build with latest changes
go build -o devtyper ./cmd/devtyper/main.go

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed! Please check the errors above."
    exit 1
fi

# Create installation directory if it doesn't exist
sudo mkdir -p /usr/local/bin

# Move binary to installation directory
sudo mv devtyper /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/devtyper

echo "DevTyper installed successfully!"
echo "Try running: devtyper --help"
