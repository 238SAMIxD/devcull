package cleaner

import (
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

func (g *GradleCleaner) IsInstalled() bool {
	p, err := g.getCachePath()
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func (g *GradleCleaner) EstimateReclaimable() (int64, error) {
	p, err := g.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(p)
}

func (g *GradleCleaner) Clean(dryRun bool) (int64, error) {
	p, err := g.getCachePath()
	if err != nil {
		return 0, err
	}
	return cleanDirs([]string{p}, dryRun)
}
