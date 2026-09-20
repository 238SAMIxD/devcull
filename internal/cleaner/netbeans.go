package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type NetBeansCleaner struct{}

func (c *NetBeansCleaner) Name() string       { return "NetBeans" }
func (c *NetBeansCleaner) Category() Category { return CategoryIDE }

func (c *NetBeansCleaner) Aliases() []string { return nil }

func (c *NetBeansCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return []string{filepath.Join(home, "Library", "Caches", "NetBeans")}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return nil
		}
		return []string{filepath.Join(localAppData, "NetBeans", "Cache")}
	default:
		return []string{filepath.Join(home, ".cache", "netbeans")}
	}
}

func (c *NetBeansCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *NetBeansCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *NetBeansCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
