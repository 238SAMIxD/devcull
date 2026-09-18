package cleaner

import (
	"os"
	"path/filepath"
)

type EclipseCleaner struct{}

func (c *EclipseCleaner) Name() string       { return "Eclipse" }
func (c *EclipseCleaner) Category() Category { return CategoryIDE }

func (c *EclipseCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, ".eclipse"),
		filepath.Join(home, ".p2", "pool", "plugins"), 
	}
}

func (c *EclipseCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func (c *EclipseCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths()), nil
}

func (c *EclipseCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}