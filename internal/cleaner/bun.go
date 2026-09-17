package cleaner

import (
	"os"
	"os/exec"
	"strings"
)

type BunCleaner struct{}

func (b *BunCleaner) Name() string {
	return "Bun"
}

func (b *BunCleaner) Category() string {
	return "Runtimes"
}

func (b *BunCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("bun"); err != nil {
		return false
	}
	if err := exec.Command("bun", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (b *BunCleaner) EstimateReclaimable() (int64, error) {
	if !b.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	out, err := exec.Command("bun", "pm", "cache").Output()
	if err != nil {
		return 0, nil
	}

	cachePath := strings.TrimSpace(string(out))
	if cachePath == "" {
		return 0, nil
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return 0, nil
	}

	return dirSize(cachePath)
}

func (b *BunCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := b.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun || reclaimable == 0 {
		return reclaimable, nil
	}

	out, err := exec.Command("bun", "pm", "cache").Output()
	if err == nil {
		cachePath := strings.TrimSpace(string(out))
		if cachePath != "" {
			_ = os.RemoveAll(cachePath)
		}
	}

	return reclaimable, nil
}