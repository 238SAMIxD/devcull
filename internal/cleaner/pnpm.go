package cleaner

import (
	"context"
	"time"

	"os/exec"
	"strings"
)

type PnpmCleaner struct{}

func (p *PnpmCleaner) Name() string {
	return "pnpm"
}

func (p *PnpmCleaner) Category() Category {
	return CategoryNode
}

func (p *PnpmCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("pnpm"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "pnpm", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PnpmCleaner) getCachePath() string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "pnpm", "store", "path").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (p *PnpmCleaner) EstimateReclaimable() (int64, error) {
	cachePath := p.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (p *PnpmCleaner) Clean(dryRun bool) (int64, error) {
	before, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "pnpm", "store", "prune").Run(); err != nil {
		return 0, err
	}

	after, err := dirSize(p.getCachePath())
	if err != nil && after == 0 {
		after = before
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, err
}
