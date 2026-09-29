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

func (b *BrewCleaner) Aliases() []string { return nil }

func (b *BrewCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("brew"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if err := exec.CommandContext(ctx, "brew", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BrewCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
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
	cachePath, err := b.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	if !isSafeToDelete(cachePath) {
		return 0, fmt.Errorf("Homebrew cache path rejected by safety guard: %s", cachePath)
	}
	return dirSize(ctx, cachePath)
}

func (b *BrewCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	cachePath, err := b.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	if !isSafeToDelete(cachePath) {
		return 0, fmt.Errorf("Homebrew cache path rejected by safety guard: %s", cachePath)
	}
	return cleanDirs(ctx, []string{cachePath}, dryRun)
}
