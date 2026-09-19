package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type VSCodeCleaner struct{}

func (c *VSCodeCleaner) Name() string       { return "VS Code" }
func (c *VSCodeCleaner) Category() Category { return CategoryIDE }

func (c *VSCodeCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var baseDir string
	switch runtime.GOOS {
	case "darwin":
		baseDir = filepath.Join(home, "Library", "Application Support", "Code")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return nil
		}
		baseDir = filepath.Join(appData, "Code")
	default:
		baseDir = filepath.Join(home, ".config", "Code")
	}

	return []string{
		filepath.Join(baseDir, "Cache"),
		filepath.Join(baseDir, "CachedData"),
		filepath.Join(baseDir, "CachedExtensionVSIXs"),
	}
}

func (c *VSCodeCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if _, err := os.Stat(filepath.Dir(p)); err == nil {
			return true
		}
	}
	return false
}

func (c *VSCodeCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths())
}

func (c *VSCodeCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}
