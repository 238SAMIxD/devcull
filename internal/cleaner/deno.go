package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
)

type DenoCleaner struct{}

func (d *DenoCleaner) Name() string {
	return "Deno"
}

func (d *DenoCleaner) Category() Category {
	return CategoryNode
}

func (d *DenoCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("deno"); err != nil {
		return false
	}
	if err := exec.Command("deno", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (d *DenoCleaner) getCachePath() string {
	if custom := os.Getenv("DENO_DIR"); custom != "" {
		if filepath.IsAbs(custom) {
			return custom
		}
		if abs, err := filepath.Abs(custom); err == nil {
			return abs
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if d, err := os.UserCacheDir(); err == nil {
		return filepath.Join(d, "deno")
	}

	return filepath.Join(home, ".cache", "deno")
}

func (d *DenoCleaner) EstimateReclaimable() (int64, error) {
	path := d.getCachePath()
	if path == "" {
		return 0, nil
	}
	return dirSize(path)
}

func (d *DenoCleaner) Clean(dryRun bool) (int64, error) {
	before, err := d.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return before, nil
	}

	path := d.getCachePath()
	if path != "" {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			return 0, err
		}
	}

	after, err := dirSize(path)
	if err != nil {
		return 0, err
	}
	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
