package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type CondaCleaner struct{}

func (c *CondaCleaner) Name() string       { return "Conda" }
func (c *CondaCleaner) Category() Category { return CategoryPython }

func (c *CondaCleaner) Aliases() []string {
	return []string{"anaconda", "miniconda", "miniforge", "conda", "anaconda3", "miniconda3", "miniforge3", "mamba", "micromamba", "mambaforge"}
}

func (c *CondaCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	paths := []string{
		filepath.Join(home, ".conda", "pkgs"),
		filepath.Join(home, "anaconda3", "pkgs"),
		filepath.Join(home, "miniconda3", "pkgs"),
		filepath.Join(home, "miniforge3", "pkgs"),
		filepath.Join(home, "mambaforge", "pkgs"),
		filepath.Join(home, "micromamba", "pkgs"),
	}

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "conda", "pkgs"))
		}
	}

	return paths
}

func (c *CondaCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *CondaCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	var validPaths []string
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			validPaths = append(validPaths, p)
		}
	}
	if len(validPaths) == 0 {
		return 0, nil
	}
	return dirsSize(ctx, validPaths)
}

func (c *CondaCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	var validPaths []string
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			validPaths = append(validPaths, p)
		}
	}
	if len(validPaths) == 0 {
		return 0, nil
	}
	return cleanDirs(ctx, validPaths, dryRun)
}
