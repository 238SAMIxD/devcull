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

func (d *DenoCleaner) Category() string {
	return "Runtimes"
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
		return custom
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
	if !d.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	path := d.getCachePath()
	if path == "" {
		return 0, nil
	}
	return dirSize(path)
}

func (d *DenoCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := d.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	path := d.getCachePath()
	if path != "" {
		_ = os.RemoveAll(path)
	}

	return reclaimable, nil
}