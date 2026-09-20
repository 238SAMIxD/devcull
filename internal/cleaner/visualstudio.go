package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type VisualStudioCleaner struct{}

func (c *VisualStudioCleaner) Name() string       { return "Visual Studio" }
func (c *VisualStudioCleaner) Category() Category { return CategoryIDE }

func (c *VisualStudioCleaner) getPaths() []string {
	if runtime.GOOS != "windows" {
		return nil
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return nil
	}
	base := filepath.Join(localAppData, "Microsoft", "VisualStudio")

	componentCaches, _ := filepath.Glob(filepath.Join(base, "*", "ComponentModelCache"))
	designerCaches, _ := filepath.Glob(filepath.Join(base, "*", "Designer", "Cache"))

	return append(componentCaches, designerCaches...)
}

func (c *VisualStudioCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *VisualStudioCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *VisualStudioCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
