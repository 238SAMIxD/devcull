package cleaner

import (
	"context"
	"fmt"
	"time"

	"os/exec"
	"strings"
)

type BrewCleaner struct{}

func (b *BrewCleaner) Name() string {
	return "Homebrew"
}

func (b *BrewCleaner) Category() Category {
	return CategorySystem
}

func (b *BrewCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("brew"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "brew", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BrewCleaner) getCachePath() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "brew", "--cache").Output()
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return p, nil
}

func (b *BrewCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := b.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (b *BrewCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := b.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "brew", "cleanup").Run(); err != nil {
		return 0, err
	}

	cachePath, err := b.getCachePath()
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
