package plugin

import (
	"context"

	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

func TestSubprocessCleaner(t *testing.T) {
	tempDir := t.TempDir()

	mockScript := filepath.Join(tempDir, "mock.sh")
	scriptContent := `#!/bin/bash
case "$1" in
  "installed") echo '{"installed": true}' ;;
  "estimate")  echo '{"reclaimable_bytes": 1024, "error": ""}' ;;
  "clean")     echo '{"reclaimed_bytes": 1024, "error": ""}' ;;
esac
`
	os.WriteFile(mockScript, []byte(scriptContent), 0755)

	p := &SubprocessCleaner{
		manifest: Manifest{
			Name:       "Mock",
			Category:   cleaner.Category("Test"),
			Entrypoint: []string{"bash", "mock.sh"},
			WorkingDir: tempDir,
		},
	}

	if !p.IsInstalled(context.Background()) {
		t.Error("Expected IsInstalled to be true")
	}

	size, err := p.EstimateReclaimable(context.Background())
	if err != nil {
		t.Errorf("Unexpected error from estimate: %v", err)
	}
	if size != 1024 {
		t.Errorf("Expected estimate size 1024, got %d", size)
	}

	reclaimed, err := p.Clean(context.Background(), false)
	if err != nil {
		t.Errorf("Unexpected error from clean: %v", err)
	}
	if reclaimed != 1024 {
		t.Errorf("Expected reclaimed size 1024, got %d", reclaimed)
	}
}

func TestSubprocessCleanerCancellation(t *testing.T) {
	tempDir := t.TempDir()

	mockScript := filepath.Join(tempDir, "mock_sleep.sh")
	scriptContent := `#!/bin/bash
trap "" SIGTERM SIGINT
sleep 10
`
	os.WriteFile(mockScript, []byte(scriptContent), 0755)

	p := &SubprocessCleaner{
		manifest: Manifest{
			Name:       "MockSleep",
			Category:   cleaner.Category("Test"),
			Entrypoint: []string{"bash", "mock_sleep.sh"},
			WorkingDir: tempDir,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := p.EstimateReclaimable(ctx)
	duration := time.Since(start)

	if err == nil {
		t.Error("Expected error from cancelled context, got nil")
	}
	if duration > 5*time.Second {
		t.Errorf("Cancellation failed or leaked, took %v", duration)
	}
}
