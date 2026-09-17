package cleaner

import (
	"os"
	"path/filepath"
)

type MavenCleaner struct{}

func (m *MavenCleaner) Name() string { return "Maven" }
func (m *MavenCleaner) Category() Category { return CategoryJava }

func (m *MavenCleaner) getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".m2", "repository")
}

func (m *MavenCleaner) IsInstalled() bool {
	info, err := os.Stat(m.getCachePath())
	return err == nil && info.IsDir()
}

func (m *MavenCleaner) EstimateReclaimable() (int64, error) {
	return dirSize(m.getCachePath())
}

func (m *MavenCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := m.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	_ = os.RemoveAll(m.getCachePath())
	return reclaimable, nil
}