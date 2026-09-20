package cleaner

import (
	"os"
	"path/filepath"
)

type NvmCleaner struct{}

func (n *NvmCleaner) Name() string {
	return "nvm"
}

func (n *NvmCleaner) Category() Category {
	return CategoryNode
}

func (n *NvmCleaner) IsInstalled() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	nvmDir := filepath.Join(home, ".nvm")
	info, err := os.Stat(nvmDir)
	return err == nil && info.IsDir()
}

func (n *NvmCleaner) getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".nvm", ".cache"), nil
}

func (n *NvmCleaner) EstimateReclaimable() (int64, error) {
	cachePath, err := n.getCachePath()
	if err != nil {
		return 0, err
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return 0, nil
	}

	return dirSize(cachePath)
}

func (n *NvmCleaner) Clean(dryRun bool) (int64, error) {
	before, err := n.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	cachePath, err := n.getCachePath()
	if err != nil {
		return 0, err
	}
	if cachePath != "" {
		if err := os.RemoveAll(cachePath); err != nil && !os.IsNotExist(err) {
			return 0, err
		}
	}

	after, err := dirSize(cachePath)
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
