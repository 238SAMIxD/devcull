package cleaner

import (
	"context"
	"time"

	"os/exec"
	"strings"
)

type DotnetCleaner struct{}

func (d *DotnetCleaner) Name() string {
	return "Dotnet"
}

func (d *DotnetCleaner) Category() Category {
	return CategoryCSharp
}

func (d *DotnetCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "dotnet", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (d *DotnetCleaner) getCachePaths(ctx context.Context) []string {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dotnet", "nuget", "locals", "all", "--list").Output()
	if err != nil {
		return nil
	}

	var paths []string
	lines := strings.Split(string(out), "\n")

	prefixes := []string{"http-cache:", "global-packages:", "temp:", "plugins-cache:"}

	for _, line := range lines {
		for _, prefix := range prefixes {
			if idx := strings.Index(line, prefix); idx != -1 {
				path := strings.TrimSpace(line[idx+len(prefix):])
				if path != "" {
					paths = append(paths, path)
				}
				break
			}
		}
	}
	return paths
}

func (d *DotnetCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, d.getCachePaths(ctx))
}

func (d *DotnetCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := d.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun || before == 0 {
		return before, nil
	}

	ctx2, cancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel2()
	if err := exec.CommandContext(ctx2, "dotnet", "nuget", "locals", "all", "--clear").Run(); err != nil {
		return 0, err
	}

	after, err := dirsSize(ctx, d.getCachePaths(ctx))
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
