package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type VSCodeCleaner struct{}

func (c *VSCodeCleaner) Name() string       { return "VS Code" }
func (c *VSCodeCleaner) Category() Category { return CategoryIDE }

func (c *VSCodeCleaner) Aliases() []string { return []string{"vscode", "code"} }

func (c *VSCodeCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var baseDir string
	switch runtime.GOOS {
	case "darwin":
		baseDir = filepath.Join(home, "Library", "Application Support", "Code")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return nil
		}
		baseDir = filepath.Join(appData, "Code")
	default:
		baseDir = filepath.Join(home, ".config", "Code")
	}

	return []string{
		filepath.Join(baseDir, "Cache"),
		filepath.Join(baseDir, "CachedData"),
		filepath.Join(baseDir, "CachedExtensionVSIXs"),
	}
}

func (c *VSCodeCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *VSCodeCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *VSCodeCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
