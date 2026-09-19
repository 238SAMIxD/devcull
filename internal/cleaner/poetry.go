package cleaner

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type PoetryCleaner struct{}

func (p *PoetryCleaner) Name() string {
	return "Poetry"
}

func (p *PoetryCleaner) Category() Category {
	return CategoryPython
}

func (p *PoetryCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("poetry"); err != nil {
		return false
	}
	if err := exec.Command("poetry", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (p *PoetryCleaner) getCachePaths() []string {
	out, err := exec.Command("poetry", "config", "cache-dir").Output()
	if err != nil {
		return nil
	}
	basePath := strings.TrimSpace(string(out))
	if basePath == "" {
		return nil
	}
	return []string{
		filepath.Join(basePath, "cache"),
		filepath.Join(basePath, "artifacts"),
	}
}

func (p *PoetryCleaner) EstimateReclaimable() (int64, error) {
	paths := p.getCachePaths()
	if len(paths) == 0 {
		return 0, nil
	}
	return dirsSize(paths)
}

func (p *PoetryCleaner) Clean(dryRun bool) (int64, error) {
	paths := p.getCachePaths()
	return cleanDirs(paths, dryRun)
}
