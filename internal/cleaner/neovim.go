package cleaner

import (
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
		return []string{filepath.Join(os.Getenv("LOCALAPPDATA"), "nvim-data", "swap")}
	}
	return []string{
		filepath.Join(home, ".cache", "nvim"),
		filepath.Join(home, ".local", "state", "nvim", "swap"),
	}
}

func (c *NeovimCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *NeovimCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths()), nil
}

func (c *NeovimCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}
