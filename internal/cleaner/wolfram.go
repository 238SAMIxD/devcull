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
		paths = append(paths, filepath.Join(home, "Library", "Caches", "Wolfram", "Mathematica"))
	case "linux":
		paths = append(paths, filepath.Join(home, ".cache", "Wolfram", "Mathematica"))
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "Wolfram Research", "Mathematica", "Cache"))
		}
	}
	return paths
}

func (c *WolframCleaner) IsInstalled(ctx context.Context) bool {
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
