package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Executor interface {
	Call(ctx context.Context, command string, arg ...string) error
	CallString(ctx context.Context, command string, arg ...string) (string, error)
}

type Cmd interface {
	CombinedOutput() ([]byte, error)
}

func NewExecutor() Executor {
	executor := executor{
		useSSHTunnel: os.Getenv("USE_SSH_TUNNEL") != "",
	}

	err := executor.Call(context.Background(), "echo", "executor initialized")
	if err != nil {
		panic(fmt.Errorf("failed to initialize executor: %w", err))
	}

	return &executor
}

type executor struct {
	useSSHTunnel bool
}

func (e *executor) CallString(ctx context.Context, command string, args ...string) (string, error) {

	var cmd *exec.Cmd

	if e.useSSHTunnel {
		sshHost := os.Getenv("SSH_HOST")
		sshPassword := os.Getenv("SSH_PASSWORD")
		originalCommand := mergeCommand(command, args)

		cmd = exec.CommandContext(ctx, "sshpass", "-p", sshPassword, "ssh",
			"-o", "StrictHostKeyChecking=no", fmt.Sprintf("root@%s", sshHost), originalCommand)
	} else {
		cmd = exec.CommandContext(ctx, command, args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logOutputLines(ctx, slog.LevelError, command, output)
		return string(output), fmt.Errorf("call failed: %s %v error=%w", command, args, err)
	}
	return string(output), nil
}

func mergeCommand(command string, args []string) string {
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			escaped := strings.ReplaceAll(arg, `"`, `\"`)
			arg = `"` + escaped + `"`
		}
		command += ` ` + arg
	}
	return command
}

func (e *executor) Call(ctx context.Context, command string, arg ...string) error {

	output, err := e.CallString(ctx, command, arg...)
	if err != nil {
		return fmt.Errorf("%w\nOutput: %s", err, output)
	}
	logOutputLines(ctx, slog.LevelInfo, command, []byte(output))
	return nil
}

func logOutputLines(ctx context.Context, level slog.Level, command string, output []byte) {

	if len(output) == 0 {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		slog.Log(ctx, level, "output", "line", line, "command", command)
	}
}
