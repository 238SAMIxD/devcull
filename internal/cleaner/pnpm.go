package cleaner

import (
	"context"
	"fmt"
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

func (p *PnpmCleaner) getCachePath() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "pnpm", "store", "path").Output()
	if err != nil {
		return "", err
	}
	pathStr := strings.TrimSpace(string(out))
	if pathStr == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return pathStr, nil
}

func (p *PnpmCleaner) EstimateReclaimable() (int64, error) {
	cachePath, err := p.getCachePath()
	if err != nil {
		return 0, err
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

	cachePath, err := p.getCachePath()
	if err != nil {
		return 0, err
	}
	after, err := dirSize(cachePath)
	if err != nil && after == 0 {
		after = before
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, err
}
