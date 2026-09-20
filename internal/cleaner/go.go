package cleaner

import (
	"context"
	"fmt"
	"time"

	"os/exec"
	"strings"
)

type GoCleaner struct{}

func (g *GoCleaner) Name() string {
	return "Go"
}

func (g *GoCleaner) Category() Category {
	return CategoryGo
}

func (g *GoCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("go"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "go", "version").Run(); err != nil {
		return false
	}
	return true
}

func (g *GoCleaner) getCachePaths() ([]string, error) {
	var paths []string
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if out, err := exec.CommandContext(ctx, "go", "env", "GOCACHE").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			paths = append(paths, p)
		}
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel2()
	if out, err := exec.CommandContext(ctx2, "go", "env", "GOMODCACHE").Output(); err == nil {
		if p := strings.TrimSpace(string(out)); p != "" {
			paths = append(paths, p)
		}
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no go cache paths found")
	}
	return paths, nil
}

func (g *GoCleaner) EstimateReclaimable() (int64, error) {
	paths, err := g.getCachePaths()
	if err != nil {
		return 0, err
	}
	return dirsSize(paths)
}

func (g *GoCleaner) Clean(dryRun bool) (int64, error) {
	before, err := g.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, "go", "clean", "-cache", "-modcache").Run(); err != nil {
		return 0, err
	}

	paths, err := g.getCachePaths()
	if err != nil {
		return 0, err
	}
	after, err := dirsSize(paths)
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
