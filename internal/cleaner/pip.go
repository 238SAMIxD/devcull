package cleaner

import (
	"os/exec"
	"strings"
)

type PipCleaner struct{}

func (p *PipCleaner) Name() string {
	return "pip"
}

func (p *PipCleaner) Category() string {
	return "Package Managers"
}

func (p *PipCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("pip"); err != nil {
		return false
	}
	if err := exec.Command("pip", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PipCleaner) EstimateReclaimable() (int64, error) {
	if !p.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("pip", "cache", "dir").Output()
	if err != nil {
		return 0, err
	}

	return dirSize(strings.TrimSpace(string(out)))
}

func (p *PipCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("pip", "cache", "purge").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}