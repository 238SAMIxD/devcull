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

func (p *PnpmCleaner) Aliases() []string { return nil }

func (p *PnpmCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("pnpm"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "pnpm", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PnpmCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
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

func (p *PnpmCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := p.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (p *PnpmCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := p.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, "pnpm", "store", "prune").Run(); err != nil {
		return 0, err
	}

	cachePath, err := p.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	after, err := dirSize(ctx, cachePath)
	if err != nil && after == 0 {
		after = before
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, err
}
