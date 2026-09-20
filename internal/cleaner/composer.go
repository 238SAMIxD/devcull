package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type ComposerCleaner struct{}

func (c *ComposerCleaner) Name() string       { return "Composer" }
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

func (c *ComposerCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *ComposerCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getCachePaths())
}

func (c *ComposerCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getCachePaths(), dryRun)
}
