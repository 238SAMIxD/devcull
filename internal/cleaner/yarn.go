package cleaner

import (
	"context"
	"fmt"
	"time"

	"os/exec"
	"strings"
)

type YarnCleaner struct{}

func (y *YarnCleaner) Name() string {
	return "Yarn"
}

func (y *YarnCleaner) Category() Category {
	return CategoryNode
}

func (y *YarnCleaner) Aliases() []string { return nil }

func (y *YarnCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("yarn"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "yarn", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (y *YarnCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "yarn", "cache", "dir").Output()
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return p, nil
}

func (y *YarnCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := y.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (y *YarnCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := y.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, "yarn", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	cachePath, err := y.getCachePath(ctx)
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
