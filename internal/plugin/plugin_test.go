package plugin

import (
	"os"
	"path/filepath"
	"testing"

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

	if !p.IsInstalled() {
		t.Error("Expected IsInstalled to be true")
	}

	size, err := p.EstimateReclaimable()
	if err != nil {
		t.Errorf("Unexpected error from estimate: %v", err)
	}
	if size != 1024 {
		t.Errorf("Expected estimate size 1024, got %d", size)
	}

	reclaimed, err := p.Clean(false)
	if err != nil {
		t.Errorf("Unexpected error from clean: %v", err)
	}
	if reclaimed != 1024 {
		t.Errorf("Expected reclaimed size 1024, got %d", reclaimed)
	}
}