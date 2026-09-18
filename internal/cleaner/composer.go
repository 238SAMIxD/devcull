package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type ComposerCleaner struct{}

func (c *ComposerCleaner) Name() string { return "Composer" }
func (c *ComposerCleaner) Category() Category { return CategoryPHP }

func (c *ComposerCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "Composer"))
		}
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, "Composer", "cache"))
		}
	} else {
		paths = append(paths, filepath.Join(home, ".composer", "cache"))
		paths = append(paths, filepath.Join(home, ".cache", "composer"))
	}

	return paths
}

func (c *ComposerCleaner) IsInstalled() bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *ComposerCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, p := range c.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (c *ComposerCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := c.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range c.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}