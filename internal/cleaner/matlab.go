package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type MatlabCleaner struct{}

func (c *MatlabCleaner) Name() string       { return "MATLAB" }
func (c *MatlabCleaner) Category() Category { return CategoryMath }
func (c *MatlabCleaner) Aliases() []string  { return []string{"matlab"} }

func (c *MatlabCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string
	switch runtime.GOOS {
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "MathWorks", "MatlabRuntimeCache"))
		}
	default: // darwin and linux
		paths = append(paths, filepath.Join(home, ".mathworks"))
	}
	return paths
}

func (c *MatlabCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *MatlabCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *MatlabCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
