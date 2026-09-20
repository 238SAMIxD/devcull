package cleaner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"os/exec"
	"strings"
)

type PipCleaner struct {
	cmdName string
	once    sync.Once
}

func (p *PipCleaner) Name() string {
	return "pip"
}

func (p *PipCleaner) Category() Category {
	return CategoryPython
}

func (p *PipCleaner) Aliases() []string { return nil }

func (p *PipCleaner) getCmd() string {
	p.once.Do(func() {
		if _, err := exec.LookPath("pip"); err == nil {
			p.cmdName = "pip"
		} else if _, err := exec.LookPath("pip3"); err == nil {
			p.cmdName = "pip3"
		}
	})
	return p.cmdName
}

func (p *PipCleaner) IsInstalled(ctx context.Context) bool {
	cmd := p.getCmd()
	if cmd == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, cmd, "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PipCleaner) getCachePath(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, p.getCmd(), "cache", "dir").Output()
	if err != nil {
		return "", err
	}
	pathStr := strings.TrimSpace(string(out))
	if pathStr == "" {
		return "", fmt.Errorf("empty cache path returned")
	}
	return pathStr, nil
}

func (p *PipCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	cachePath, err := p.getCachePath(ctx)
	if err != nil {
		return 0, err
	}
	return dirSize(ctx, cachePath)
}

func (p *PipCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := p.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, p.getCmd(), "cache", "purge").Run(); err != nil {
		return 0, err
	}

	cachePath, err := p.getCachePath(ctx)
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
