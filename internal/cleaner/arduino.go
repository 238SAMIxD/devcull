package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

type ArduinoCleaner struct{}

func (c *ArduinoCleaner) Name() string {
	return "Arduino"
}

func (c *ArduinoCleaner) Category() Category {
	return CategoryIDE
}

func (c *ArduinoCleaner) Aliases() []string {
	return []string{"arduino-ide", "arduino-cli", "arduino-lab"}
}

func (c *ArduinoCleaner) getPaths() []string {
	var paths []string
	home, err := os.UserHomeDir()
	if err != nil {
		return paths
	}

	switch runtime.GOOS {
	case "darwin":
		paths = append(paths,
			filepath.Join(home, "Library", "Caches", "arduino"),
			filepath.Join(home, "Library", "Caches", "cc.arduino.IDE2"),
			filepath.Join(home, "Library", "Application Support", "arduino-ide", "Cache"),
			filepath.Join(home, "Library", "Application Support", "arduino-ide", "CachedData"),
		)
	case "linux":
		paths = append(paths,
			filepath.Join(home, ".cache", "arduino"),
			filepath.Join(home, ".config", "arduino-ide", "Cache"),
		)
	case "windows":
		if localApp := os.Getenv("LOCALAPPDATA"); localApp != "" {
			paths = append(paths, filepath.Join(localApp, "arduino"))
		}
		if temp := os.Getenv("TEMP"); temp != "" {
			paths = append(paths, filepath.Join(temp, "arduino"))
		}
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, "arduino-ide", "Cache"))
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

func (c *ArduinoCleaner) IsInstalled(ctx context.Context) bool {
	return len(c.getPaths()) > 0
}

func (c *ArduinoCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	return dirsSize(ctx, c.getPaths())
}

func (c *ArduinoCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	return cleanDirs(ctx, c.getPaths(), dryRun)
}
