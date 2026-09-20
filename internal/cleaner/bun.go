package cleaner

import (
	"context"
	"fmt"
	"time"

	"os"
	"os/exec"
	"path/filepath"
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

func (b *BunCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "pm", "cache").Output()
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", err
		}
		if os.PathSeparator == '\\' {
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData != "" {
				return filepath.Join(localAppData, "bun", "install", "cache"), nil
			}
			return filepath.Join(home, "AppData", "Local", "bun", "install", "cache"), nil
		}
		return filepath.Join(home, ".bun", "install", "cache"), nil
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", fmt.Errorf("empty cache path returned")
		}
		if os.PathSeparator == '\\' {
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData != "" {
				return filepath.Join(localAppData, "bun", "install", "cache"), nil
			}
			return filepath.Join(home, "AppData", "Local", "bun", "install", "cache"), nil
		}
		return filepath.Join(home, ".bun", "install", "cache"), nil
	}
	return p, nil
}

func (b *BunCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := b.getCachePath(ctx)
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

	cachePath, err := b.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	if !isSafeToDelete(cachePath) {
		return 0, nil
	}
	if cachePath != "" {
		if err := removeAll(ctx, cachePath); err != nil && !os.IsNotExist(err) {
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
