package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type DartCleaner struct{}

func (p *DartCleaner) Name() string       { return "Dart" }
func (p *DartCleaner) Category() Category { return CategoryFlutter }

func (p *DartCleaner) Aliases() []string { return nil }

func (p *DartCleaner) getCachePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return []string{
				filepath.Join(localAppData, "Pub", "Cache", "hosted"),
				filepath.Join(localAppData, "Pub", "Cache", "git"),
			}
		}
		return []string{
			filepath.Join(home, "AppData", "Local", "Pub", "Cache", "hosted"),
			filepath.Join(home, "AppData", "Local", "Pub", "Cache", "git"),
		}
	}

	return []string{
		filepath.Join(home, ".pub-cache", "hosted"),
		filepath.Join(home, ".pub-cache", "git"),
	}
}

func (p *DartCleaner) IsInstalled(ctx context.Context) bool {
	for _, path := range p.getCachePaths() {
		if info, err := os.Stat(filepath.Dir(path)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (p *DartCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, p.getCachePaths())
}

func (p *DartCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, p.getCachePaths(), dryRun)
}
