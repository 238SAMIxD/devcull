package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type CocoaPodsCleaner struct{}

func (c *CocoaPodsCleaner) Name() string { return "CocoaPods" }
func (c *CocoaPodsCleaner) Category() Category { return CategoryApple }

func (c *CocoaPodsCleaner) getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Caches", "CocoaPods")
	}
	
	return filepath.Join(home, ".cocoapods")
}

func (c *CocoaPodsCleaner) IsInstalled() bool {
	info, err := os.Stat(c.getCachePath())
	return err == nil && info.IsDir()
}

func (c *CocoaPodsCleaner) EstimateReclaimable() (int64, error) {
	return dirSize(c.getCachePath())
}

func (c *CocoaPodsCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := c.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	_ = os.RemoveAll(c.getCachePath())
	return reclaimable, nil
}