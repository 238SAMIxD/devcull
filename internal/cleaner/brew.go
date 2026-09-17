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

func (b *BrewCleaner) getCachePath() string {
	out, err := exec.Command("brew", "--cache").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (b *BrewCleaner) EstimateReclaimable() (int64, error) {
	if !b.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	cachePath := b.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (b *BrewCleaner) Clean(dryRun bool) (int64, error) {
	before, err := b.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	if err := exec.Command("brew", "cleanup").Run(); err != nil {
		return 0, err
	}

	after, _ := dirSize(b.getCachePath())
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}