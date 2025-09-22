package wifi

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/d-rk/checkin-system/pkg/cmd"
)

//go:embed wifi-manager.sh
var wifiManagerScript string

type Service interface {
	ListNetworks(ctx context.Context) ([]string, error)
	AddNetwork(ctx context.Context, ssid, password string) error
	RemoveNetwork(ctx context.Context, ssid string) error
	GetWifiMode(ctx context.Context) (bool, error)
	ToggleWifiMode(ctx context.Context) (bool, error)
}

type service struct {
	executor   cmd.Executor
	scriptPath string
}

func NewService() Service {

	scriptPath, err := writeScriptToTemp()
	if err != nil {
		panic(err)
	}

	return &service{
		executor:   cmd.NewExecutor(),
		scriptPath: scriptPath,
	}
}

func (s *service) ListNetworks(ctx context.Context) ([]string, error) {
	output, err := s.executeScriptString(ctx, "list")
	if err != nil {
		return nil, err
	}

	output = strings.TrimSpace(output)

	if output == "No networks configured" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

func (s *service) AddNetwork(ctx context.Context, ssid, password string) error {
	return s.executeScript(ctx, "add", ssid, password)
}

func (s *service) RemoveNetwork(ctx context.Context, ssid string) error {
	return s.executeScript(ctx, "remove", ssid)
}

func (s *service) GetWifiMode(ctx context.Context) (bool, error) {

	output, err := s.executeScriptString(ctx, "mode")
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)
	return output == "hotspot", nil
}

func (s *service) ToggleWifiMode(ctx context.Context) (bool, error) {
	if err := s.executeScript(ctx, "toggle-mode"); err != nil {
		return false, err
	}
	return s.GetWifiMode(ctx)
}

func writeScriptToTemp() (string, error) {
	tmpDir := os.TempDir()
	scriptPath := filepath.Join(tmpDir, "wifi-manager.sh")

	err := os.WriteFile( //nolint:gosec // we want the script to be executable
		scriptPath,
		[]byte(wifiManagerScript),
		0755,
	)
	if err != nil {
		return "", fmt.Errorf("failed to write wifi script to temp file: %w", err)
	}

	return scriptPath, nil
}

func (s *service) executeScriptString(ctx context.Context, args ...string) (string, error) {
	return s.executor.CallString(ctx, s.scriptPath, args...)
}

func (s *service) executeScript(ctx context.Context, args ...string) error {
	return s.executor.Call(ctx, s.scriptPath, args...)
}
