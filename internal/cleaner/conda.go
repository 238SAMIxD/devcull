package cleaner

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type CondaCleaner struct {
	cmdName string
	once    sync.Once
	paths   []string
	pathsMu sync.Mutex
}

func (c *CondaCleaner) Name() string       { return "Conda" }
func (c *CondaCleaner) Category() Category { return CategoryPython }

func (c *CondaCleaner) Aliases() []string {
	return []string{"anaconda", "miniconda", "miniforge", "conda", "anaconda3", "miniconda3", "miniforge3", "mamba", "micromamba", "mambaforge"}
}

func (c *CondaCleaner) getCmd() string {
	c.once.Do(func() {
		for _, name := range []string{"conda", "mamba", "micromamba"} {
			if path, err := exec.LookPath(name); err == nil {
				c.cmdName = path
				return
			}
		}
		
		home, err := os.UserHomeDir()
		if err == nil {
			var baseDirs []string
			baseDirs = append(baseDirs,
				filepath.Join(home, "miniconda3"),
				filepath.Join(home, "anaconda3"),
				filepath.Join(home, "miniforge3"),
				filepath.Join(home, "mambaforge"),
				filepath.Join(home, "micromamba"),
			)
			for _, base := range baseDirs {
				bin := filepath.Join(base, "bin", "conda")
				if runtime.GOOS == "windows" {
					bin = filepath.Join(base, "Scripts", "conda.exe")
				}
				if _, err := os.Stat(bin); err == nil {
					c.cmdName = bin
					return
				}
			}
		}
	})
	return c.cmdName
}

type condaInfo struct {
	PkgsDirs []string `json:"pkgs_dirs"`
}

func (c *CondaCleaner) getPaths(ctx context.Context) []string {
	c.pathsMu.Lock()
	if c.paths != nil {
		c.pathsMu.Unlock()
		return c.paths
	}
	c.pathsMu.Unlock()

	cmdName := c.getCmd()
	if cmdName == "" {
		return nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx2, cmdName, "info", "--json").Output()
	if err != nil {
		return nil
	}

	var info condaInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil
	}

	var validPaths []string
	for _, p := range info.PkgsDirs {
		validPaths = append(validPaths, p)
	}

	c.pathsMu.Lock()
	c.paths = validPaths
	c.pathsMu.Unlock()

	return validPaths
}

func (c *CondaCleaner) IsInstalled(ctx context.Context) bool {
	if c.getCmd() == "" {
		return false
	}
	
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return exec.CommandContext(ctx2, c.getCmd(), "--version").Run() == nil
}

func (c *CondaCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	var validPaths []string
	for _, p := range c.getPaths(ctx) {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			validPaths = append(validPaths, p)
		}
	}
	if len(validPaths) == 0 {
		return 0, nil
	}
	return dirsSize(ctx, validPaths)
}

func (c *CondaCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	before, err := c.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	if dryRun || before == 0 {
		return before, nil
	}

	cmdName := c.getCmd()
	if cmdName == "" {
		return 0, nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx2, cmdName, "clean", "--all", "--yes")
	if err := cmd.Run(); err != nil {
		return 0, err
	}

	after, err := c.EstimateReclaimable(ctx)
	if err != nil {
		return 0, err
	}

	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
