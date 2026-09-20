package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type FvmCleaner struct{}

func (f *FvmCleaner) Name() string       { return "FVM" }
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

func (f *FvmCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range f.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (f *FvmCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, f.getCachePaths())
}

func (f *FvmCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, f.getCachePaths(), dryRun)
}
