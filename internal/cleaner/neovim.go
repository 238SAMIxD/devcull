package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type NeovimCleaner struct{}

func (c *NeovimCleaner) Name() string       { return "Neovim" }
func (c *NeovimCleaner) Category() Category { return CategoryIDE }

func (c *NeovimCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		return nil
	}
	return []string{
		filepath.Join(home, ".cache", "nvim"),
	}
}

func (c *NeovimCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *NeovimCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *NeovimCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
