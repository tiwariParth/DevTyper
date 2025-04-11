package monitor

import "strings"

type CommandType int

const (
	Generic CommandType = iota
	Docker
	Kubernetes
	NPM
	Go
)

var commonCommands = map[string]struct {
	Type        CommandType
	Description string
}{
	"docker pull":      {Docker, "Pulling Docker image"},
	"docker build":     {Docker, "Building Docker image"},
	"kubectl apply":    {Kubernetes, "Applying Kubernetes manifests"},
	"eksctl create":    {Kubernetes, "Creating EKS cluster"},
	"npm install":      {NPM, "Installing NPM packages"},
	"yarn install":     {NPM, "Installing Yarn packages"},
	"go mod download":  {Go, "Downloading Go dependencies"},
	"create-next-app": {NPM, "Creating Next.js application"},
}

var interactiveCommands = map[string]struct{
    Title       string
    Args        []string
    Prompt      string
}{
    "create-next-app": {
        Title: "Create Next.js App",
        Args: []string{
            "--name <project-name>",
            "--typescript",
            "--tailwind",
            "--eslint",
            "--app",
            "--src-dir",
            "--import-alias",
        },
        Prompt: "To skip interactive mode, provide arguments like:\ndevtyper npx create-next-app@latest my-app --typescript --tailwind",
    },
    "npm init": {
        Title: "NPM Init",
        Args: []string{"-y"},
        Prompt: "Use 'devtyper npm init -y' to skip interactive mode",
    },
}

func DetectCommand(cmd string) (CommandType, string, bool) {
    for pattern, info := range commonCommands {
        if strings.HasPrefix(cmd, pattern) {
            // Check if this is an interactive command
            interactive, _, _ := DetectInteractiveCommand(cmd)
            return info.Type, info.Description, interactive
        }
    }
    return Generic, "Running command", false
}

func DetectInteractiveCommand(cmd string) (bool, string, string) {
    for pattern, info := range interactiveCommands {
        if strings.Contains(cmd, pattern) {
            return true, info.Title, info.Prompt
        }
    }
    return false, "", ""
}
