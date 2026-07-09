package setupcommand

import (
	"context"
	"fmt"
	"strings"

	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/parser"

	docker "gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/setupCommand/docker"
	shell "gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/setupCommand/shell"

	"github.com/moby/moby/client"
)

type SetupCommand interface {
	Execute() error
	Cleanup() error
}

func NewSetupCommand(dto *parser.SetupCommandDTO, cli *client.Client, baseDir string, ctx context.Context) (SetupCommand, error) {
	handler := docker.NewDockerHandler(cli, baseDir, ctx)

	switch strings.ToLower(dto.Type) {
	case "shell":
		return shell.NewShellSetupCommand(dto, baseDir, ctx)
	case "docker":
		return docker.NewDockerSetupCommand(dto, handler)
	default:
		return nil, fmt.Errorf("unsupported setup command type: %s", dto.Type)
	}
}

func NewSetupCommands(dtos []*parser.SetupCommandDTO, cli *client.Client, baseDir string, ctx context.Context) ([]SetupCommand, error) {
	var setupCommands []SetupCommand
	for _, dto := range dtos {
		cmd, err := NewSetupCommand(dto, cli, baseDir, ctx)
		if err != nil {
			return nil, fmt.Errorf("error creating setup command: %w", err)
		}
		if cmd != nil {
			setupCommands = append(setupCommands, cmd)
		}
	}
	return setupCommands, nil
}
