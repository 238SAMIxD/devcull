package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type DartCleaner struct{}

func (p *DartCleaner) Name() string { return "Dart" }
func (p *DartCleaner) Category() Category { return CategoryFlutter }

func (p *DartCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return []string{filepath.Join(localAppData, "Pub", "Cache")}
		}
		return []string{filepath.Join(home, "AppData", "Local", "Pub", "Cache")}
	}

	return []string{filepath.Join(home, ".pub-cache")}
}

func (p *DartCleaner) IsInstalled() bool {
	for _, path := range p.getCachePaths() {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (p *DartCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, path := range p.getCachePaths() {
		size, _ := dirSize(path)
		total += size
	}
	return total, nil
}

func (p *DartCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := p.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, path := range p.getCachePaths() {
		_ = os.RemoveAll(path)
	}
	return reclaimable, nil
}