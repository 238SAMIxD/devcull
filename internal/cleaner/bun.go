package cleaner

import (
	"context"
	"fmt"
	"time"

	"os"
	"os/exec"
	"strings"
)

type BunCleaner struct{}

func (b *BunCleaner) Name() string {
	return "Bun"
}

func (b *BunCleaner) Category() Category {
	return CategoryNode
}

func (b *BunCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("bun"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "bun", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BunCleaner) getCachePath() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "pm", "cache").Output()
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return p, nil
}

func (b *BunCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := b.getCachePath()
	if err != nil {
		return 0, err
	}

	if !isSafeToDelete(cachePath) {
		return 0, nil
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return 0, nil
	}

	return dirSize(ctx, cachePath)
}

func (b *BunCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := b.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun || before == 0 {
		return before, nil
	}

	cachePath, err := b.getCachePath()
	if err != nil {
		return 0, err
	}
	if !isSafeToDelete(cachePath) {
		return 0, nil
	}
	if cachePath != "" {
		if err := os.RemoveAll(cachePath); err != nil && !os.IsNotExist(err) {
			return 0, err
		}
	}

	after, err := dirSize(ctx, cachePath)
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
