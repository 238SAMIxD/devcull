package cleaner

import (
	"context"
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

func (b *BunCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("bun"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "bun", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BunCleaner) getCachePath() string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "pm", "cache").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (b *BunCleaner) EstimateReclaimable() (int64, error) {
	cachePath := b.getCachePath()
	if cachePath == "" {
		return 0, nil
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return 0, nil
	}

	return dirSize(cachePath)
}

func (b *BunCleaner) Clean(dryRun bool) (int64, error) {
	before, err := b.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun || before == 0 {
		return before, nil
	}

	cachePath := b.getCachePath()
	if cachePath != "" {
		if err := os.RemoveAll(cachePath); err != nil && !os.IsNotExist(err) {
			return 0, err
		}
	}

	after, err := dirSize(cachePath)
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
