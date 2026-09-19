package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromDir(t *testing.T) {
	tempDir := t.TempDir()

	validDir := filepath.Join(tempDir, "valid-plugin")
	os.MkdirAll(validDir, 0755)
	validJSON := `{"name": "Valid Plugin", "category": "System & DevOps", "entrypoint": ["bash", "test.sh"]}`
	os.WriteFile(filepath.Join(validDir, "manifest.json"), []byte(validJSON), 0644)

	emptyDir := filepath.Join(tempDir, "empty-plugin")
	os.MkdirAll(emptyDir, 0755)

	corruptedDir := filepath.Join(tempDir, "corrupted-plugin")
	os.MkdirAll(corruptedDir, 0755)
	os.WriteFile(filepath.Join(corruptedDir, "manifest.json"), []byte(`{ "name": "Broken", missing_quotes }`), 0644)

	noNameDir := filepath.Join(tempDir, "no-name-plugin")
	os.MkdirAll(noNameDir, 0755)
	noNameJSON := `{"category": "System & DevOps", "entrypoint": ["bash", "test.sh"]}`
	os.WriteFile(filepath.Join(noNameDir, "manifest.json"), []byte(noNameJSON), 0644)

	noEntryDir := filepath.Join(tempDir, "no-entrypoint-plugin")
	os.MkdirAll(noEntryDir, 0755)
	noEntryJSON := `{"name": "No Entry", "category": "System & DevOps"}`
	os.WriteFile(filepath.Join(noEntryDir, "manifest.json"), []byte(noEntryJSON), 0644)

	os.WriteFile(filepath.Join(tempDir, "just-a-file.txt"), []byte("this should be ignored"), 0644)

	plugins := loadFromDir(tempDir)

	if len(plugins) != 1 {
		t.Fatalf("Expected exactly 1 valid plugin to survive the edge cases, got %d", len(plugins))
	}

	if plugins[0].Name() != "Valid Plugin" {
		t.Errorf("Expected the valid plugin to be loaded, got '%s'", plugins[0].Name())
	}
}
