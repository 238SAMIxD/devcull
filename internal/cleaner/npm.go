package cleaner

import (
	"context"
	"fmt"
	"time"

	"os/exec"
	"strings"
)

type NpmCleaner struct{}

func (n *NpmCleaner) Name() string {
	return "npm"
}

func (n *NpmCleaner) Category() Category {
	return CategoryNode
}

func (n *NpmCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("npm"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "npm", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (n *NpmCleaner) getCachePath() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "npm", "config", "get", "cache").Output()
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return p, nil
}

func (n *NpmCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := n.getCachePath()
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (n *NpmCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := n.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "npm", "cache", "clean", "--force").Run(); err != nil {
		return 0, err
	}

	cachePath, err := n.getCachePath()
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
