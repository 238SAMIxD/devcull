package cleaner

import (
	"os"
	"path/filepath"
)

type CcacheCleaner struct{}

func (c *CcacheCleaner) Name() string { return "Ccache" }
func (c *CcacheCleaner) Category() Category { return CategoryCpp }

func (c *CcacheCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	
	paths := []string{
		filepath.Join(home, ".ccache"),
		filepath.Join(home, ".cache", "ccache"),
	}
	
	paths = append(paths, filepath.Join(home, "Library", "Caches", "ccache"))
	
	return paths
}

func (c *CcacheCleaner) IsInstalled() bool {
	for _, p := range c.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *CcacheCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, p := range c.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (c *CcacheCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := c.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range c.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}