package cleaner

import (
	"os"
	"os/exec"
	"strings"
)

type PoetryCleaner struct{}

func (p *PoetryCleaner) Name() string {
	return "Poetry"
}

func (p *PoetryCleaner) Category() string {
	return "Package Managers"
}

func (p *PoetryCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("poetry"); err != nil {
		return false
	}
	if err := exec.Command("poetry", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PoetryCleaner) EstimateReclaimable() (int64, error) {
	if !p.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("poetry", "config", "cache-dir").Output()
	if err != nil {
		return 0, err
	}

	cachePath := strings.TrimSpace(string(out))
	if cachePath == "" {
		return 0, nil
	}

	return dirSize(cachePath)
}

func (p *PoetryCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	out, err := exec.Command("poetry", "config", "cache-dir").Output()
	if err == nil {
		cachePath := strings.TrimSpace(string(out))
		if cachePath != "" {
			_ = os.RemoveAll(cachePath)
		}
	}

	return reclaimable, nil
}