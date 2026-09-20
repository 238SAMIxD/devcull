package cleaner

import (
	"context"
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

func (u *UvCleaner) getCachePath() string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "uv", "cache", "dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (u *UvCleaner) EstimateReclaimable() (int64, error) {
	cachePath := u.getCachePath()
	if cachePath == "" {
		return 0, nil
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

	after, err := dirSize(u.getCachePath())
	if err != nil && after == 0 {
		after = before
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, err
}
