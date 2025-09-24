package wifi

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/d-rk/checkin-system/pkg/cmd"
)

//go:embed wifi-manager.sh
var wifiManagerScript string

const outputSeparator = "==========="

type Status struct {
	State     string
	IPAddress *string
	SSID      *string
	Mode      string
}

type Service interface {
	ListNetworks(ctx context.Context) ([]string, error)
	AddNetwork(ctx context.Context, ssid, password string) error
	RemoveNetwork(ctx context.Context, ssid string) error
	GetStatus(ctx context.Context) (Status, error)
	ToggleWifiMode(ctx context.Context) error
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

	s := service{
		executor:   cmd.NewExecutor(),
		scriptPath: scriptPath,
	}

	status, err := s.GetStatus(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get initial wifi status: %w", err))
	}

	if status.Mode == "client" && (strings.ToLower(status.State) == "down" || status.IPAddress == nil) {
		slog.Warn("wifi in client mode but no IP address assigned")
		err = s.ToggleWifiMode(context.Background())
		if err != nil {
			panic(fmt.Errorf("failed to toggle wifi mode: %w", err))
		}
		slog.Info("wifi mode updated")
	}

	return &s
}

func (s *service) ListNetworks(ctx context.Context) ([]string, error) {
	output, err := s.executeScriptString(ctx, "list")
	if err != nil {
		return nil, err
	}

	output = filterOutput(output)

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

func (s *service) GetStatus(ctx context.Context) (Status, error) {
	output, err := s.executeScriptString(ctx, "status")
	if err != nil {
		return Status{}, err
	}

	output = filterOutput(output)

	return Status{
		State:     strings.ToLower(extractField(output, "State")),
		IPAddress: extractFieldPtr(output, "IP Address"),
		SSID:      extractFieldPtr(output, "SSID"),
		Mode:      strings.ToLower(extractField(output, "Mode")),
	}, nil
}

func extractFieldPtr(output string, key string) *string {
	// Create regex pattern to match "key: value" format
	pattern := fmt.Sprintf(`(?m)^%s:\s*(.+)$`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 { //nolint:mnd // regex has to parts
		value := strings.TrimSpace(matches[1])
		// Return nil for certain "empty" values
		if value == "" || strings.HasPrefix(value, "None") {
			return nil
		}
		return &value
	}

	return nil
}

func extractField(output string, key string) string {
	value := extractFieldPtr(output, key)
	if value == nil {
		return ""
	}
	return *value
}

func (s *service) ToggleWifiMode(ctx context.Context) error {
	return s.executeScript(ctx, "toggle-mode")
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

func filterOutput(output string) string {
	// Check if output contains separator and extract part after it
	if idx := strings.Index(output, outputSeparator); idx != -1 {
		output = output[idx+len(outputSeparator):]
	} else {
		output = ""
	}

	return strings.TrimSpace(output)
}
