package cleaner

import (
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

	base := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "VisualStudio")

	componentCaches, _ := filepath.Glob(filepath.Join(base, "*", "ComponentModelCache"))
	designerCaches, _ := filepath.Glob(filepath.Join(base, "*", "Designer", "Cache"))

	return append(componentCaches, designerCaches...)
}

func (c *VisualStudioCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *VisualStudioCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths())
}

func (c *VisualStudioCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}
