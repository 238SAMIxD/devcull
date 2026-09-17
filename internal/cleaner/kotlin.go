package cleaner

import (
	"os"
	"path/filepath"
)

type KotlinCleaner struct{}

func (k *KotlinCleaner) Name() string { return "Kotlin" }
func (k *KotlinCleaner) Category() Category { return CategoryJava }

func (k *KotlinCleaner) getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".konan")
}

func (k *KotlinCleaner) IsInstalled() bool {
	info, err := os.Stat(k.getCachePath())
	return err == nil && info.IsDir()
}

func (k *KotlinCleaner) EstimateReclaimable() (int64, error) {
	return dirSize(k.getCachePath())
}

func (k *KotlinCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := k.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	_ = os.RemoveAll(k.getCachePath())
	return reclaimable, nil
}