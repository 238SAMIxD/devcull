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

func (y *YarnCleaner) EstimateReclaimable() (int64, error) {
	if !y.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("yarn", "cache", "dir").Output()
	if err != nil {
		return 0, err
	}

	return dirSize(strings.TrimSpace(string(out)))
}

func (y *YarnCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := y.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("yarn", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}