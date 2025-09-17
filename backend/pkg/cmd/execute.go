package cmd

import (
	"bufio"
	"context"
	"log/slog"
	"os/exec"
	"sync"
)

type Executor interface {
	Call(ctx context.Context, command string, arg ...string) error
}

type Cmd interface {
	CombinedOutput() ([]byte, error)
}

func NewExecutor() Executor {
	return &executor{}
}

type executor struct{}

func (executor) Call(ctx context.Context, command string, arg ...string) error {
	cmd := exec.CommandContext(ctx, command, arg...)

	stdout, stdOutErr := cmd.StdoutPipe()
	if stdOutErr != nil {
		return stdOutErr
	}

	stderr, stdErrErr := cmd.StderrPipe()
	if stdErrErr != nil {
		return stdErrErr
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	logger := slog.Default().With(slog.Group("exec"))
	var wg sync.WaitGroup
	wg.Add(2) //nolint:mnd // two goroutines

	// Handle stdout in a goroutine
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logger.Debug(command, "stdout", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			logger.Error("error reading stdout", "error", err)
		}
	}()

	// Handle stderr in a goroutine
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logger.Error(command, "stderr", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			logger.Error("error reading stderr", "error", err)
		}
	}()

	err := cmd.Wait()
	wg.Wait()
	return err
}
