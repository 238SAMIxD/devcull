package cleaner

import (
	"os"
	"path/filepath"
)

type AndroidCleaner struct{}

func (a *AndroidCleaner) Name() string       { return "Android" }
func (a *AndroidCleaner) Category() Category { return CategoryJava }

func (a *AndroidCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".android", "cache"),
		filepath.Join(home, ".android", "build-cache"),
	}
}

func (a *AndroidCleaner) IsInstalled() bool {
	for _, p := range a.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (a *AndroidCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(a.getCachePaths())
}

func (a *AndroidCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(a.getCachePaths(), dryRun)
}
