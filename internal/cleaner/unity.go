package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type UnityCleaner struct{}

func (u *UnityCleaner) Name() string       { return "Unity" }
func (u *UnityCleaner) Category() Category { return CategoryCSharp }

func (u *UnityCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string

	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			filepath.Join(home, "Library", "Unity", "cache", "packages"),
			filepath.Join(home, "Library", "Unity", "Asset Store-5.x"),
			filepath.Join(home, "Library", "Caches", "Unity", "gi_cache"),
			filepath.Join(home, "Library", "Application Support", "UnityHub", "Downloads"), // Hub Installers
		}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			if !filepath.IsAbs(localAppData) {
				if abs, err := filepath.Abs(localAppData); err == nil {
					localAppData = abs
				} else {
					localAppData = ""
				}
			}
		}

		appData := os.Getenv("APPDATA")
		if appData != "" {
			if !filepath.IsAbs(appData) {
				if abs, err := filepath.Abs(appData); err == nil {
					appData = abs
				} else {
					appData = ""
				}
			}
		}

		if localAppData != "" {
			paths = append(paths,
				filepath.Join(localAppData, "Unity", "cache", "packages"),
				filepath.Join(localAppData, "Unity", "cache", "gi_cache"),
			)
		}
		if appData != "" {
			paths = append(paths,
				filepath.Join(appData, "Unity", "Asset Store-5.x"),
				filepath.Join(appData, "UnityHub", "Downloads"),
			)
		}
	default:
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil
		}
		paths = []string{
			filepath.Join(configDir, "unity3d", "cache", "packages"),
			filepath.Join(home, ".local", "share", "unity3d", "Asset Store-5.x"),
			filepath.Join(configDir, "unity3d", "cache", "gi_cache"),
			filepath.Join(configDir, "UnityHub", "Downloads"),
		}
	}

	return paths
}

func (u *UnityCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range u.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (u *UnityCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, u.getCachePaths())
}

func (u *UnityCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, u.getCachePaths(), dryRun)
}
