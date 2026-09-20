package cleaner

import (
	"os"
	"path/filepath"
)

type KotlinCleaner struct{}

func (k *KotlinCleaner) Name() string       { return "Kotlin" }
func (k *KotlinCleaner) Category() Category { return CategoryJava }

func (k *KotlinCleaner) getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".konan"), nil
}

func (k *KotlinCleaner) IsInstalled() bool {
	p, err := k.getCachePath()
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func (k *KotlinCleaner) EstimateReclaimable() (int64, error) {
	p, err := k.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(p)
}

func (k *KotlinCleaner) Clean(dryRun bool) (int64, error) {
	p, err := k.getCachePath()
	if err != nil {
		return 0, err
	}
	return cleanDirs([]string{p}, dryRun)
}
