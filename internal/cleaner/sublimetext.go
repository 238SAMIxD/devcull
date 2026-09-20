package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type SublimeCleaner struct{}

func (c *SublimeCleaner) Name() string       { return "Sublime Text" }
func (c *SublimeCleaner) Category() Category { return CategoryIDE }

func (c *SublimeCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return []string{
			filepath.Join(home, "Library", "Caches", "Sublime Text"),
			filepath.Join(home, "Library", "Caches", "Sublime Text 3"),
		}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return nil
		}
		return []string{filepath.Join(localAppData, "Sublime Text", "Cache")}
	default:
		return []string{filepath.Join(home, ".cache", "sublime-text")}
	}
}

func (c *SublimeCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *SublimeCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *SublimeCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
