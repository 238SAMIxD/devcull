package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type WolframCleaner struct{}

func (c *WolframCleaner) Name() string       { return "Wolfram Mathematica" }
func (c *WolframCleaner) Category() Category { return CategoryMath }
func (c *WolframCleaner) Aliases() []string  { return []string{"wolfram", "mathematica", "wolfram-alpha", "wolfram-desktop", "wolframalpha"} }

func (c *WolframCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		if matches, err := filepath.Glob(filepath.Join(home, "Library", "Mathematica", "FrontEnd", "*_Caches")); err == nil {
			paths = append(paths, matches...)
		}
	case "linux":
		if matches, err := filepath.Glob(filepath.Join(home, ".Mathematica", "FrontEnd", "*_Caches")); err == nil {
			paths = append(paths, matches...)
		}
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			if matches, err := filepath.Glob(filepath.Join(localAppData, "Mathematica", "FrontEnd", "* Caches")); err == nil {
				paths = append(paths, matches...)
			}
		}
	}
	return paths
}

func (c *WolframCleaner) IsInstalled(ctx context.Context) bool {
	home, err := os.UserHomeDir()
	if err == nil {
		var baseDir string
		switch runtime.GOOS {
		case "darwin":
			baseDir = filepath.Join(home, "Library", "Mathematica")
		case "linux":
			baseDir = filepath.Join(home, ".Mathematica")
		case "windows":
			if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
				baseDir = filepath.Join(localAppData, "Mathematica")
			}
		}
		if baseDir != "" {
			if info, err := os.Stat(baseDir); err == nil && info.IsDir() {
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

func (c *WolframCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *WolframCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
