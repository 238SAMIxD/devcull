package cleaner

import (
	"context"
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

func (c *CargoCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("cargo"); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "cargo", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (c *CargoCleaner) getCachePaths() []string {
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
		return nil
	}

	return []string{
		filepath.Join(cargoHome, "registry", "cache"),
		filepath.Join(cargoHome, "registry", "src"),
		filepath.Join(cargoHome, "git"),
	}
}

func (c *CargoCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getCachePaths())
}

func (c *CargoCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getCachePaths(), dryRun)
}
