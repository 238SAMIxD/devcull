package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type CocoaPodsCleaner struct{}

func (c *CocoaPodsCleaner) Name() string       { return "CocoaPods" }
func (c *CocoaPodsCleaner) Category() Category { return CategoryApple }

func (c *CocoaPodsCleaner) Aliases() []string { return nil }

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

func (c *CocoaPodsCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *CocoaPodsCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getCachePaths())
}

func (c *CocoaPodsCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getCachePaths(), dryRun)
}
