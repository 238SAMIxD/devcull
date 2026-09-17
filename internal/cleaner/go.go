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

func (g *GoCleaner) getCachePaths() []string {
	var paths []string
	if out, err := exec.Command("go", "env", "GOCACHE").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			paths = append(paths, p)
		}
	}
	if out, err := exec.Command("go", "env", "GOMODCACHE").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

func (g *GoCleaner) EstimateReclaimable() (int64, error) {
	if !g.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	var total int64
	for _, p := range g.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (g *GoCleaner) Clean(dryRun bool) (int64, error) {
	before, err := g.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	if err := exec.Command("go", "clean", "-cache", "-modcache").Run(); err != nil {
		return 0, err
	}

	var after int64
	for _, p := range g.getCachePaths() {
		size, _ := dirSize(p)
		after += size
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}