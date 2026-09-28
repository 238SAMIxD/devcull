package cleaner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type WingetCleaner struct{}

func (c *WingetCleaner) Name() string {
	return "WinGet"
}

func (c *WingetCleaner) Category() Category {
	return CategorySystem
}

func (c *WingetCleaner) Aliases() []string {
	return []string{"winget", "windows package manager"}
}

func (c *WingetCleaner) IsInstalled(ctx context.Context) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	if _, err := exec.LookPath("winget"); err == nil {
		return true
	}

	paths := c.getCachePaths()
	return len(paths) > 0
}

func (c *WingetCleaner) getCachePaths() []string {
	var paths []string

	if tempDir := os.Getenv("TEMP"); tempDir != "" {
		paths = append(paths, filepath.Join(tempDir, "WinGet"))
	}

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		paths = append(paths, filepath.Join(localAppData, "Packages", "Microsoft.DesktopAppInstaller_8wekyb3d8bbwe", "TempState", "WinGet"))
	}

	var existingPaths []string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			existingPaths = append(existingPaths, p)
		}
	}

	return existingPaths
}

func (c *WingetCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	paths := c.getCachePaths()
	if len(paths) == 0 {
		return 0, nil
	}
	return dirsSize(ctx, paths)
}

func (c *WingetCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	paths := c.getCachePaths()
	if len(paths) == 0 {
		return 0, nil
	}
	return cleanDirs(ctx, paths, dryRun)
}
