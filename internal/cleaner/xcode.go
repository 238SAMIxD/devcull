package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type XcodeCleaner struct{}

func (x *XcodeCleaner) Name() string       { return "Xcode" }
func (x *XcodeCleaner) Category() Category { return CategoryApple }

func (x *XcodeCleaner) Aliases() []string { return nil }

func (x *XcodeCleaner) getCachePaths() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData"),
		filepath.Join(home, "Library", "Caches", "com.apple.dt.Xcode"),
	}
}

func (x *XcodeCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range x.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (x *XcodeCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, x.getCachePaths())
}

func (x *XcodeCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, x.getCachePaths(), dryRun)
}
