package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type PhpbrewCleaner struct{}

func (p *PhpbrewCleaner) Name() string       { return "phpbrew" }
func (p *PhpbrewCleaner) Category() Category { return CategoryPHP }

func (p *PhpbrewCleaner) Aliases() []string { return nil }

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

func (p *PhpbrewCleaner) IsInstalled(ctx context.Context) bool {
	for _, path := range p.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(path)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (p *PhpbrewCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, p.getCachePaths())
}

func (p *PhpbrewCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, p.getCachePaths(), dryRun)
}
