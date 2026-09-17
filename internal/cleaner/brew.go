package cleaner

import (
	"os/exec"
	"strings"
)

type BrewCleaner struct{}

func (b *BrewCleaner) Name() string {
	return "Homebrew"
}

func (b *BrewCleaner) Category() string {
	return "Package Managers"
}

func (b *BrewCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("brew"); err != nil {
		return false
	}
	if err := exec.Command("brew", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BrewCleaner) EstimateReclaimable() (int64, error) {
	if !b.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("brew", "--cache").Output()
	if err != nil {
		return 0, err
	}

	cachePath := strings.TrimSpace(string(out))
	return dirSize(cachePath)
}

func (b *BrewCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := b.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("brew", "cleanup").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}