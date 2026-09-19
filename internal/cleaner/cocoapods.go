package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type CocoaPodsCleaner struct{}

func (c *CocoaPodsCleaner) Name() string       { return "CocoaPods" }
func (c *CocoaPodsCleaner) Category() Category { return CategoryApple }

func (c *CocoaPodsCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	if runtime.GOOS == "darwin" {
		return []string{filepath.Join(home, "Library", "Caches", "CocoaPods")}
	}

	return nil
}

func (c *CocoaPodsCleaner) IsInstalled() bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *CocoaPodsCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getCachePaths())
}

func (c *CocoaPodsCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getCachePaths(), dryRun)
}
