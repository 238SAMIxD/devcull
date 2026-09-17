package cleaner

import (
	"os/exec"
	"strings"
)

type PipCleaner struct {
	cmdName string
}

func (p *PipCleaner) Name() string {
	return "pip"
}

func (p *PipCleaner) Category() Category {
	return CategoryPython
}

func (p *PipCleaner) getCmd() string {
	if p.cmdName != "" {
		return p.cmdName
	}
	if _, err := exec.LookPath("pip"); err == nil {
		p.cmdName = "pip"
		return p.cmdName
	}
	if _, err := exec.LookPath("pip3"); err == nil {
		p.cmdName = "pip3"
		return p.cmdName
	}
	return ""
}

func (p *PipCleaner) IsInstalled() bool {
	cmd := p.getCmd()
	if cmd == "" {
		return false
	}
	if err := exec.Command(cmd, "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PipCleaner) getCachePath() string {
	out, err := exec.Command(p.getCmd(), "cache", "dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (p *PipCleaner) EstimateReclaimable() (int64, error) {
	cachePath := p.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (p *PipCleaner) Clean(dryRun bool) (int64, error) {
	before, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	if err := exec.Command(p.getCmd(), "cache", "purge").Run(); err != nil {
		return 0, err
	}

	after, _ := dirSize(p.getCachePath())
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}