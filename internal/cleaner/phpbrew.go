package cleaner

import (
	"os"
	"path/filepath"
)

type PhpbrewCleaner struct{}

func (p *PhpbrewCleaner) Name() string { return "phpbrew" }
func (p *PhpbrewCleaner) Category() Category { return CategoryPHP }

func (p *PhpbrewCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, ".phpbrew", "distfiles"),
		filepath.Join(home, ".phpbrew", "build"),
	}
}

func (p *PhpbrewCleaner) IsInstalled() bool {
	for _, path := range p.getCachePaths() {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (p *PhpbrewCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, path := range p.getCachePaths() {
		size, _ := dirSize(path)
		total += size
	}
	return total, nil
}

func (p *PhpbrewCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := p.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, path := range p.getCachePaths() {
		_ = os.RemoveAll(path)
	}
	return reclaimable, nil
}