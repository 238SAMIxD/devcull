package cleaner

import (
	"os/exec"
	"strings"
)

type NpmCleaner struct{}

func (n *NpmCleaner) Name() string {
	return "npm"
}

func (n *NpmCleaner) Category() string {
	return "Package Managers"
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

func (n *NpmCleaner) EstimateReclaimable() (int64, error) {
	if !n.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("npm", "config", "get", "cache").Output()
	if err != nil {
		return 0, err
	}

	cachePath := strings.TrimSpace(string(out))
	return dirSize(cachePath)
}

func (n *NpmCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := n.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("npm", "cache", "clean", "--force").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}