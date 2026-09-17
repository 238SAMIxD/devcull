package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type SwiftPMCleaner struct{}

func (s *SwiftPMCleaner) Name() string { return "SwiftPM" }
func (s *SwiftPMCleaner) Category() Category { return CategoryApple }

func (s *SwiftPMCleaner) getCachePaths() []string {
	if runtime.GOOS != "darwin" {
		home, _ := os.UserHomeDir()
		return []string{filepath.Join(home, ".swiftpm")}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	
	return []string{
		filepath.Join(home, "Library", "Caches", "org.swift.swiftpm"),
		filepath.Join(home, "Library", "org.swift.swiftpm"),
	}
}

func (s *SwiftPMCleaner) IsInstalled() bool {
	for _, p := range s.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (s *SwiftPMCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, p := range s.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (s *SwiftPMCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := s.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range s.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}