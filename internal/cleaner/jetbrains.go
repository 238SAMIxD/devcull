package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type JetBrainsCleaner struct{}

func (c *JetBrainsCleaner) Name() string       { return "JetBrains" }
func (c *JetBrainsCleaner) Category() Category { return CategoryIDE }

func (c *JetBrainsCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return []string{filepath.Join(home, "Library", "Caches", "JetBrains")}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return nil
		}
		return []string{filepath.Join(localAppData, "JetBrains")}
	default:
		return []string{filepath.Join(home, ".cache", "JetBrains")}
	}
}

func (c *JetBrainsCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func (c *JetBrainsCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths())
}

func (c *JetBrainsCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}
