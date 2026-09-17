package cleaner

import (
	"os"
	"path/filepath"
)

type ConanCleaner struct{}

func (c *ConanCleaner) Name() string { return "Conan" }
func (c *ConanCleaner) Category() Category { return CategoryCpp }

func (c *ConanCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".conan"),
		filepath.Join(home, ".conan2"),
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
	var total int64
	for _, p := range c.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (c *ConanCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := c.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range c.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}