package cleaner

import (
	"os/exec"
	"strings"
)

type GoCleaner struct{}

func (g *GoCleaner) Name() string {
	return "Go"
}

func (g *GoCleaner) Category() string {
	return "Languages"
}

func (g *GoCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("go"); err != nil {
		return false
	}
	if err := exec.Command("go", "version").Run(); err != nil {
		return false
	}
	return true
}

func (g *GoCleaner) EstimateReclaimable() (int64, error) {
	if !g.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	var total int64

	if out, err := exec.Command("go", "env", "GOCACHE").Output(); err == nil {
		size, _ := dirSize(strings.TrimSpace(string(out)))
		total += size
	}

	if out, err := exec.Command("go", "env", "GOMODCACHE").Output(); err == nil {
		size, _ := dirSize(strings.TrimSpace(string(out)))
		total += size
	}

	return total, nil
}

func (g *GoCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := g.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("go", "clean", "-cache", "-modcache").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}