package cleaner

import (
	"context"
	"errors"
	"time"

	"os/exec"
	"path/filepath"
	"strings"
)

type PoetryCleaner struct{}

func (p *PoetryCleaner) Name() string {
	return "Poetry"
}

func (p *PoetryCleaner) Category() Category {
	return CategoryPython
}

func (p *PoetryCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("poetry"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "poetry", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PoetryCleaner) getCachePaths(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "poetry", "config", "cache-dir").Output()
	if err != nil {
		return nil, err
	}
	basePath := strings.TrimSpace(string(out))
	if basePath == "" {
		return nil, errors.New("empty poetry cache dir")
	}
	return []string{
		filepath.Join(basePath, "cache"),
		filepath.Join(basePath, "artifacts"),
	}, nil
}

func (p *PoetryCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	paths, err := p.getCachePaths(ctx)
	if err != nil {
		return 0, err
	}
	if len(paths) == 0 {
		return 0, nil
	}
	return dirsSize(ctx, paths)
}

func (p *PoetryCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	paths, err := p.getCachePaths(ctx)
	if err != nil {
		return 0, err
	}
	return cleanDirs(ctx, paths, dryRun)
}
