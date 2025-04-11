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

var interactiveCommands = map[string]struct {
	Description string
	ArgExample  string
}{
	"create-next-app": {
		Description: "Next.js app creation requires interactive input",
		ArgExample:  "npx create-next-app@latest my-app --typescript --tailwind --eslint --app",
	},
	"npm init": {
		Description: "NPM init requires interactive input",
		ArgExample:  "npm init -y",
	},
}

func DetectCommand(cmd string) (CommandType, string, bool, string) {
	// Check for interactive commands first
	for pattern, info := range interactiveCommands {
		if strings.Contains(cmd, pattern) {
			return Generic, info.Description, true, info.ArgExample
		}
	}

	// Check regular commands
	for pattern, info := range commonCommands {
		if strings.HasPrefix(cmd, pattern) {
			return info.Type, info.Description, false, ""
		}
	}
	return Generic, "Running command", false, ""
}
