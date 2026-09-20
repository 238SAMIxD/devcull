package cleaner

import (
	"context"

	"os"
	"path/filepath"
)

type MavenCleaner struct{}

func (m *MavenCleaner) Name() string       { return "Maven" }
func (m *MavenCleaner) Category() Category { return CategoryJava }

func (m *MavenCleaner) getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".m2", "repository"), nil
}

func (m *MavenCleaner) IsInstalled(ctx context.Context) bool {
	p, err := m.getCachePath()
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func (m *MavenCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	p, err := m.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, p)
}

func (m *MavenCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	p, err := m.getCachePath()
	if err != nil {
		return 0, err
	}
	return cleanDirs(ctx, []string{p}, dryRun)
}
