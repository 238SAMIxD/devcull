package cleaner

import (
	"context"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "deno", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (d *DenoCleaner) getCachePaths() []string {
	var base string
	if custom := os.Getenv("DENO_DIR"); custom != "" {
		if filepath.IsAbs(custom) {
			base = custom
		} else if abs, err := filepath.Abs(custom); err == nil {
			base = abs
		}
	}

	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		if cacheDir, err := os.UserCacheDir(); err == nil {
			base = filepath.Join(cacheDir, "deno")
		} else {
			base = filepath.Join(home, ".cache", "deno")
		}
	}

	if base == "" {
		return nil
	}

	base = filepath.Clean(base)
	if base == "/" || base == "." || base == "\\" {
		return nil
	}
	if home, err := os.UserHomeDir(); err == nil && base == filepath.Clean(home) {
		return nil
	}
	vol := filepath.VolumeName(base)
	if base == vol+"\\" || base == vol+"/" || base == vol {
		return nil
	}

	return []string{
		filepath.Join(base, "deps"),
		filepath.Join(base, "gen"),
	}
}

func (d *DenoCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(d.getCachePaths())
}

func (d *DenoCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(d.getCachePaths(), dryRun)
}
