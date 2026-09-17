package cleaner

import (
	"os/exec"
	"strings"
)

type UvCleaner struct{}

func (u *UvCleaner) Name() string {
	return "uv"
}

func (u *UvCleaner) Category() string {
	return "Package Managers"
}

func (u *UvCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("uv"); err != nil {
		return false
	}
	if err := exec.Command("uv", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (u *UvCleaner) EstimateReclaimable() (int64, error) {
	if !u.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("uv", "cache", "dir").Output()
	if err != nil {
		return 0, err
	}

	cachePath := strings.TrimSpace(string(out))
	return dirSize(cachePath)
}

func (u *UvCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := u.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	if err := exec.Command("uv", "cache", "clean").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}