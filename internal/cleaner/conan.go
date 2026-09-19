package cleaner

import (
	"os"
	"path/filepath"
)

type ConanCleaner struct{}

func (c *ConanCleaner) Name() string       { return "Conan" }
func (c *ConanCleaner) Category() Category { return CategoryCpp }

func (c *ConanCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".conan", "data"),
		filepath.Join(home, ".conan2", "p"),
	}
}

func (c *ConanCleaner) IsInstalled() bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *ConanCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getCachePaths())
}

func (c *ConanCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getCachePaths(), dryRun)
}
