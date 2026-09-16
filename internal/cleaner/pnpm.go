package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PnpmCleaner struct{}

func (p *PnpmCleaner) Name() string {
	return "pnpm"
}

func (p *PnpmCleaner) Category() string {
	return "Package Managers"
}

func (p *PnpmCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("pnpm"); err != nil {
		return false
	}
	if err := exec.Command("pnpm", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PnpmCleaner) EstimateReclaimable() (int64, error) {
	if !p.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("pnpm", "store", "path").Output()
	if err != nil {
		return 0, err
	}

	storePath := strings.TrimSpace(string(out))

	return dirSize(storePath)
}

func (p *PnpmCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := p.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("pnpm", "store", "prune").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size, err
}