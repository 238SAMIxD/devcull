package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type UnrealCleaner struct{}

func (u *UnrealCleaner) Name() string { return "Unreal Engine" }
func (u *UnrealCleaner) Category() Category { return CategoryCpp }

func (u *UnrealCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return []string{filepath.Join(home, "Library", "Application Support", "Epic", "UnrealEngine", "Common", "DerivedDataCache")}
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return []string{filepath.Join(localAppData, "UnrealEngine", "Common", "DerivedDataCache")}
		}
	default:
		if configDir, err := os.UserConfigDir(); err == nil {
			return []string{filepath.Join(configDir, "Epic", "UnrealEngine", "Common", "DerivedDataCache")}
		}
	}
	return nil
}

func (u *UnrealCleaner) IsInstalled() bool {
	for _, p := range u.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (u *UnrealCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, p := range u.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (u *UnrealCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := u.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range u.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}