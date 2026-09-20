package cleaner

import (
	"context"
	"fmt"
	"time"

	"os/exec"
	"strings"
)

type UvCleaner struct{}

func (u *UvCleaner) Name() string {
	return "uv"
}

func (u *UvCleaner) Category() Category {
	return CategoryPython
}

func (u *UvCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("uv"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "uv", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (u *UvCleaner) getCachePath() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "uv", "cache", "dir").Output()
	if err != nil {
		return "", err
	}
	pathStr := strings.TrimSpace(string(out))
	if pathStr == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return pathStr, nil
}

func (u *UvCleaner) EstimateReclaimable() (int64, error) {
	cachePath, err := u.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(cachePath)
}

func (u *UvCleaner) Clean(dryRun bool) (int64, error) {
	before, err := u.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "uv", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	cachePath, err := u.getCachePath()
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
