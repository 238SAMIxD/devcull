package cleaner

import (
	"os/exec"
	"strings"
)

type YarnCleaner struct{}

func (y *YarnCleaner) Name() string {
	return "Yarn"
}

func (y *YarnCleaner) Category() string {
	return "Package Managers"
}

func (y *YarnCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("yarn"); err != nil {
		return false
	}
	if err := exec.Command("yarn", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (y *YarnCleaner) getCachePath() string {
	out, err := exec.Command("yarn", "cache", "dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (y *YarnCleaner) EstimateReclaimable() (int64, error) {
	cachePath := y.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (y *YarnCleaner) Clean(dryRun bool) (int64, error) {
	before, err := y.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	if err := exec.Command("yarn", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	after, _ := dirSize(y.getCachePath())
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}