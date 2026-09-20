package cleaner

import (
	"context"
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

func (y *YarnCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("yarn"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "yarn", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (y *YarnCleaner) getCachePath() string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "yarn", "cache", "dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (y *YarnCleaner) EstimateReclaimable() (int64, error) {
	cachePath := y.getCachePath()
	if cachePath == "" {
		return 0, nil
	}
	return dirSize(cachePath)
}

func (y *YarnCleaner) Clean(dryRun bool) (int64, error) {
	before, err := y.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "yarn", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	after, err := dirSize(y.getCachePath())
	if err != nil && after == 0 {
		after = before
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, err
}
