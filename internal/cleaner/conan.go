package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type ConanCleaner struct{}

func (c *ConanCleaner) Name() string       { return "Conan" }
func (c *ConanCleaner) Category() Category { return CategoryCpp }

func (c *ConanCleaner) Aliases() []string { return nil }

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

func (c *ConanCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *ConanCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getCachePaths())
}

func (c *ConanCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getCachePaths(), dryRun)
}
