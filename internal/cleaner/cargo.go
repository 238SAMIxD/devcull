package cleaner

import (
	"context"
	"fmt"
	"time"

	"os"
	"os/exec"
	"path/filepath"
)

type CargoCleaner struct{}

func (c *CargoCleaner) Name() string {
	return "Cargo"
}

func (c *CargoCleaner) Category() Category {
	return CategoryRust
}

func (c *CargoCleaner) Aliases() []string { return nil }

func (c *CargoCleaner) IsInstalled(ctx context.Context) bool {
	if _, err := exec.LookPath("cargo"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "cargo", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (c *CargoCleaner) getCachePaths() ([]string, error) {
	var cargoHome string
	if home := os.Getenv("CARGO_HOME"); home != "" {
		if filepath.IsAbs(home) {
			cargoHome = home
		} else if abs, err := filepath.Abs(home); err == nil {
			cargoHome = abs
		}
	}

	if cargoHome == "" {
		if userHome, err := os.UserHomeDir(); err == nil {
			cargoHome = filepath.Join(userHome, ".cargo")
		}
	}

	if cargoHome == "" {
		return nil, fmt.Errorf("cargo home is not found")
	}

	return []string{
		filepath.Join(cargoHome, "registry", "cache"),
		filepath.Join(cargoHome, "registry", "src"),
		filepath.Join(cargoHome, "git"),
	}, nil
}

func (c *CargoCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	paths, err := c.getCachePaths()
	if err != nil {
		return 0, err
	}
	return dirsSize(ctx, paths)
}

func (c *CargoCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	paths, err := c.getCachePaths()
	if err != nil {
		return 0, err
	}
	return cleanDirs(ctx, paths, dryRun)
}
