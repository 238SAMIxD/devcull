package cleaner

import (
	"os"
	"path/filepath"
)

type NvmCleaner struct{}

func (n *NvmCleaner) Name() string {
	return "nvm"
}

func (n *NvmCleaner) Category() string {
	return "Version Managers"
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

func (n *NvmCleaner) getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".nvm", ".cache")
}

func (n *NvmCleaner) EstimateReclaimable() (int64, error) {
	cachePath := n.getCachePath()
	if cachePath == "" {
		return 0, nil
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

	cachePath := n.getCachePath()
	if cachePath != "" {
		_ = os.RemoveAll(cachePath)
	}

	after, _ := dirSize(cachePath)
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}