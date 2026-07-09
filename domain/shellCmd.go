package domain

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

func ExecuteShellCommand(cmdStr string, baseDir string, ctx context.Context) (int, error) {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return 1, fmt.Errorf("no command provided")
	}

	parts, err := shlex.Split(cmdStr)
	if err != nil {
		return 1, fmt.Errorf("error parsing command '%s': %v", cmdStr, err)
	}

	head := parts[0]
	args := parts[1:]

	cmd := exec.CommandContext(ctx, head, args...) // nolint:gosec

	if baseDir != "" {
		cmd.Dir = baseDir
	}

	output, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		return exitCode, fmt.Errorf("error executing command '%s': %v. Error details: %s", cmdStr, err, string(output))
	}

	slog.Info("Output of command executed", "command", cmdStr, "output", string(output), "exitCode", exitCode)

	return exitCode, nil
}
