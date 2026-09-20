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

func (u *UvCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("uv"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "uv", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (u *UvCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
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

func (u *UvCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := u.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (u *UvCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := u.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, "uv", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	cachePath, err := u.getCachePath(ctx)
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
