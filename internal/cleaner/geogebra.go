package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type GeoGebraCleaner struct{}

func (c *GeoGebraCleaner) Name() string       { return "GeoGebra" }
func (c *GeoGebraCleaner) Category() Category { return CategoryMath }
func (c *GeoGebraCleaner) Aliases() []string  { return []string{"geogebra"} }

func (c *GeoGebraCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths, filepath.Join(home, "Library", "Caches", "GeoGebra"))
	case "linux":
		paths = append(paths, filepath.Join(home, ".cache", "GeoGebra"))
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "GeoGebra_5.0", "Cache"))
		}
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, "GeoGebra", "Cache"))
		}
	}
	return paths
}

func (c *GeoGebraCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *GeoGebraCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *GeoGebraCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
