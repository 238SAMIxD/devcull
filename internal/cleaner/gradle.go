package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type GradleCleaner struct{}

func (g *GradleCleaner) Name() string       { return "Gradle" }
func (g *GradleCleaner) Category() Category { return CategoryJava }

func (g *GradleCleaner) getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gradle", "caches"), nil
}

func (g *GradleCleaner) IsInstalled(ctx context.Context) bool {
	p, err := g.getCachePath()
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func (g *GradleCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	p, err := g.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, p)
}

func (g *GradleCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	p, err := g.getCachePath()
	if err != nil {
		return 0, err
	}
	return cleanDirs(ctx, []string{p}, dryRun)
}
