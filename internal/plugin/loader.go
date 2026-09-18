package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

func LoadPlugins() []cleaner.Cleaner {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}

	pluginsDir := filepath.Join(configDir, "devcull", "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return nil
	}

	var plugins []cleaner.Cleaner

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pluginsDir, entry.Name())
		manifestPath := filepath.Join(pluginPath, "manifest.json")

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}

		if manifest.Name == "" || len(manifest.Entrypoint) == 0 {
			continue
		}

		manifest.WorkingDir = pluginPath
		plugins = append(plugins, &SubprocessCleaner{manifest: manifest})
	}

	return plugins
}