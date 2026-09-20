package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type UnrealCleaner struct{}

func (u *UnrealCleaner) Name() string       { return "Unreal Engine" }
func (u *UnrealCleaner) Category() Category { return CategoryCpp }

func (u *UnrealCleaner) Aliases() []string { return nil }

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

func (u *UnrealCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range u.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (u *UnrealCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, u.getCachePaths())
}

func (u *UnrealCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, u.getCachePaths(), dryRun)
}
