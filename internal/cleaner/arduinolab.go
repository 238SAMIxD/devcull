package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type ArduinoLabCleaner struct{}

func (c *ArduinoLabCleaner) Name() string {
	return "Arduino Lab"
}

func (c *ArduinoLabCleaner) Category() Category {
	return CategoryIDE
}

func (c *ArduinoLabCleaner) Aliases() []string {
	return []string{"arduino-lab", "arduino-lab-for-micropython", "micropython"}
}

func (c *ArduinoLabCleaner) getPaths() []string {
	var paths []string
	home, err := os.UserHomeDir()
	if err != nil {
		return paths
	}

	switch runtime.GOOS {
	case "darwin":
		paths = append(paths, filepath.Join(home, "Library", "Application Support", "arduino-lab-for-micropython", "Cache"))
	case "linux":
		if configDir, err := os.UserConfigDir(); err == nil {
			paths = append(paths, filepath.Join(configDir, "arduino-lab-for-micropython", "Cache"))
		}
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, "arduino-lab-for-micropython", "Cache"))
		}
	}

	var existingPaths []string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			existingPaths = append(existingPaths, p)
		}
	}
	return existingPaths
}

func (c *ArduinoLabCleaner) IsInstalled(ctx context.Context) bool {
	return len(c.getPaths()) > 0
}

func (c *ArduinoLabCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *ArduinoLabCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
