package cleaner

import (
	"os/exec"
	"strings"
)

type NpmCleaner struct{}

func (n *NpmCleaner) Name() string {
	return "npm"
}

func (n *NpmCleaner) Category() Category {
	return CategoryNode
}

func (n *NpmCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("npm"); err != nil {
		return false
	}
	if err := exec.Command("npm", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (n *NpmCleaner) getCachePath() string {
	out, err := exec.Command("npm", "config", "get", "cache").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (n *NpmCleaner) EstimateReclaimable() (int64, error) {
	cachePath := n.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (n *NpmCleaner) Clean(dryRun bool) (int64, error) {
	before, err := n.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	if err := exec.Command("npm", "cache", "clean", "--force").Run(); err != nil {
		return 0, err
	}

	after, _ := dirSize(n.getCachePath())
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
