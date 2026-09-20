package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type CcacheCleaner struct{}

func (c *CcacheCleaner) Name() string       { return "Ccache" }
func (c *CcacheCleaner) Category() Category { return CategoryCpp }

func (c *CcacheCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string
	for _, char := range "0123456789abcdef" {
		paths = append(paths, filepath.Join(home, ".ccache", string(char)))
	}
	paths = append(paths, filepath.Join(home, ".ccache", "tmp"))
	paths = append(paths, filepath.Join(home, ".cache", "ccache"))
	paths = append(paths, filepath.Join(home, "Library", "Caches", "ccache"))

	return paths
}

func (c *CcacheCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *CcacheCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getCachePaths())
}

func (c *CcacheCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getCachePaths(), dryRun)
}
