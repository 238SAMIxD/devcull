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

func (c *CondaCleaner) IsInstalled(ctx context.Context) bool {
	if c.getCmd() == "" {
		return false
	}

	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return exec.CommandContext(ctx2, c.getCmd(), "--version").Run() == nil
}

type condaCleanJSON struct {
	Tarballs struct {
		TotalSize int64 `json:"total_size"`
	} `json:"tarballs"`
	Packages struct {
		TotalSize int64 `json:"total_size"`
	} `json:"packages"`
}

func (c *CondaCleaner) getCleanSize(ctx context.Context, dryRun bool) (int64, error) {
	cmdName := c.getCmd()
	if cmdName == "" {
		return 0, nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	args := []string{"clean", "--all", "--json"}
	if dryRun {
		args = append(args, "--dry-run")
	} else {
		args = append(args, "--yes")
	}

	out, err := exec.CommandContext(ctx2, cmdName, args...).Output()
	if err != nil {
		return 0, err
	}

	var res condaCleanJSON
	if err := json.Unmarshal(out, &res); err != nil {
		return 0, err
	}

	return res.Tarballs.TotalSize + res.Packages.TotalSize, nil
}

func (c *CondaCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return c.getCleanSize(ctx, true)
}

func (c *CondaCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return c.getCleanSize(ctx, dryRun)
}
