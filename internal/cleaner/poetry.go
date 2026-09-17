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

func (p *PoetryCleaner) getCachePath() string {
	out, err := exec.Command("poetry", "config", "cache-dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (p *PoetryCleaner) EstimateReclaimable() (int64, error) {
	cachePath := p.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (p *PoetryCleaner) Clean(dryRun bool) (int64, error) {
	before, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	cachePath := p.getCachePath()
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