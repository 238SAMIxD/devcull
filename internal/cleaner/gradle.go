package cleaner

import (
	"os"
	"path/filepath"
)

type GradleCleaner struct{}

func (g *GradleCleaner) Name() string { return "Gradle" }
func (g *GradleCleaner) Category() Category { return CategoryJava }

func (g *GradleCleaner) getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".gradle", "caches")
}

func (g *GradleCleaner) IsInstalled() bool {
	info, err := os.Stat(g.getCachePath())
	return err == nil && info.IsDir()
}

func (g *GradleCleaner) EstimateReclaimable() (int64, error) {
	return dirSize(g.getCachePath())
}

func (g *GradleCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := g.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	_ = os.RemoveAll(g.getCachePath())
	return reclaimable, nil
}