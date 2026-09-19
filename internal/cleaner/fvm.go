package cleaner

import (
	"os"
	"path/filepath"
)

type FvmCleaner struct{}

func (f *FvmCleaner) Name() string { return "FVM" }
func (f *FvmCleaner) Category() Category { return CategoryFlutter }

func (f *FvmCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	
	var paths []string
	if matches, err := filepath.Glob(filepath.Join(home, "fvm", "versions", "*", "bin", "cache")); err == nil {
		paths = append(paths, matches...)
	}
	if matches, err := filepath.Glob(filepath.Join(home, ".fvm", "versions", "*", "bin", "cache")); err == nil {
		paths = append(paths, matches...)
	}
	
	return paths
}

func (f *FvmCleaner) IsInstalled() bool {
	for _, p := range f.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (f *FvmCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(f.getCachePaths()), nil
}

func (f *FvmCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(f.getCachePaths(), dryRun)
}