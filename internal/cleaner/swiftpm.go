package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type SwiftPMCleaner struct{}

func (s *SwiftPMCleaner) Name() string       { return "SwiftPM" }
func (s *SwiftPMCleaner) Category() Category { return CategoryApple }

func (s *SwiftPMCleaner) getCachePaths() []string {
	if runtime.GOOS != "darwin" {
		home, _ := os.UserHomeDir()
		return []string{filepath.Join(home, ".swiftpm", "cache")}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, "Library", "Caches", "org.swift.swiftpm"),
	}
}

func (s *SwiftPMCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range s.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (s *SwiftPMCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, s.getCachePaths())
}

func (s *SwiftPMCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, s.getCachePaths(), dryRun)
}
