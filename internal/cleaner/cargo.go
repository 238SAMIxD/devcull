package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
)

type CargoCleaner struct{}

func (c *CargoCleaner) Name() string {
	return "Cargo"
}

func (c *CargoCleaner) Category() string {
	return "Package Managers"
}

func (c *CargoCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("cargo"); err != nil {
		return false
	}
	if err := exec.Command("cargo", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (c *CargoCleaner) getCachePaths() []string {
	var cargoHome string
	if home := os.Getenv("CARGO_HOME"); home != "" {
		cargoHome = home
	} else if userHome, err := os.UserHomeDir(); err == nil {
		cargoHome = filepath.Join(userHome, ".cargo")
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
	if !c.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	paths := c.getCachePaths()
	var total int64
	for _, p := range paths {
		size, _ := dirSize(p)
		total += size
	}
	
	return total, nil
}

func (c *CargoCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := c.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	for _, p := range c.getCachePaths() {
		_ = os.RemoveAll(p)
	}

	return reclaimable, nil
}