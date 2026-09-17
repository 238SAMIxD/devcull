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

func (p *PipCleaner) Category() string {
	return "Package Managers"
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

func (p *PipCleaner) EstimateReclaimable() (int64, error) {
	if !p.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command(p.getCmd(), "cache", "dir").Output()
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

	if err := exec.Command(p.getCmd(), "cache", "purge").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}