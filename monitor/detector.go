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

func DetectCommand(cmd string) (CommandType, string) {
	for pattern, info := range commonCommands {
		if strings.HasPrefix(cmd, pattern) {
			return info.Type, info.Description
		}
	}
	return Generic, "Running command"
}
