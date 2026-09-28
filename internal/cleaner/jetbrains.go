package cleaner

import (
	"context"

	"os"
	"path/filepath"
	"runtime"
)

type JetBrainsCleaner struct{}

func (c *JetBrainsCleaner) Name() string       { return "JetBrains" }
func (c *JetBrainsCleaner) Category() Category { return CategoryIDE }

func (c *JetBrainsCleaner) Aliases() []string {
	return []string{"idea", "intellij", "pycharm", "webstorm", "goland", "rider", "clion", "phpstorm", "rubymine"}
}

func (c *JetBrainsCleaner) getPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return []string{filepath.Join(home, "Library", "Caches", "JetBrains")}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return nil
		}
		return []string{filepath.Join(localAppData, "JetBrains")}
	default:
		return []string{filepath.Join(home, ".cache", "JetBrains")}
	}
}

func (c *JetBrainsCleaner) IsInstalled(ctx context.Context) bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(filepath.Dir(p)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (c *JetBrainsCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *JetBrainsCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
