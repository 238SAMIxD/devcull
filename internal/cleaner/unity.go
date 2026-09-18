package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type UnityCleaner struct{}

func (u *UnityCleaner) Name() string { return "Unity" }
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
		appData := os.Getenv("APPDATA")
		
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
		configDir, _ := os.UserConfigDir()
		paths = []string{
			filepath.Join(configDir, "unity3d", "cache", "packages"),
			filepath.Join(home, ".local", "share", "unity3d", "Asset Store-5.x"),
			filepath.Join(configDir, "unity3d", "cache", "gi_cache"),
			filepath.Join(configDir, "UnityHub", "Downloads"),
		}
	}

	return paths
}

func (u *UnityCleaner) IsInstalled() bool {
	for _, p := range u.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (u *UnityCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(u.getCachePaths()), nil
}

func (u *UnityCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(u.getCachePaths(), dryRun)
}