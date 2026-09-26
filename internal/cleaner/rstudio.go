package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type RStudioCleaner struct{}

func (c *RStudioCleaner) Name() string       { return "R Studio" }
func (c *RStudioCleaner) Category() Category { return CategoryMath }
func (c *RStudioCleaner) Aliases() []string  { return []string{"rstudio", "r"} }

func (c *RStudioCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths, filepath.Join(home, "Library", "Caches", "RStudio"))
	case "linux":
		paths = append(paths, filepath.Join(home, ".cache", "rstudio"))
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "RStudio"))
			paths = append(paths, filepath.Join(localAppData, "RStudio-Desktop", "ctx"))
		}
	}
	return paths
}

func (c *RStudioCleaner) IsInstalled(ctx context.Context) bool {
	home, err := os.UserHomeDir()
	if err == nil {
		var configPath string
		switch runtime.GOOS {
		case "darwin":
			configPath = filepath.Join(home, "Library", "Application Support", "RStudio")
		case "linux":
			configPath = filepath.Join(home, ".local", "share", "rstudio")
		case "windows":
			if appData := os.Getenv("APPDATA"); appData != "" {
				configPath = filepath.Join(appData, "RStudio")
			}
		}
		if configPath != "" {
			if info, err := os.Stat(configPath); err == nil && info.IsDir() {
				return true
			}
		}
	}

	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *RStudioCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *RStudioCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
