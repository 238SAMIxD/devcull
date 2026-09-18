package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

func LoadPlugins() []cleaner.Cleaner {
	var plugins []cleaner.Cleaner

	if cwd, err := os.Getwd(); err == nil {
		plugins = append(plugins, loadFromDir(filepath.Join(cwd, "plugins"))...)
	}
	if configDir, err := os.UserConfigDir(); err == nil {
		plugins = append(plugins, loadFromDir(filepath.Join(configDir, "devcull", "plugins"))...)
	}

	return plugins
}

func loadFromDir(dir string) []cleaner.Cleaner {
	var plugins []cleaner.Cleaner
	entries, err := os.ReadDir(dir)
	if err != nil {
		return plugins
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(filepath.Join(pluginPath, "manifest.json"))
		if err != nil {
			continue
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil || manifest.Name == "" || len(manifest.Entrypoint) == 0 {
				continue
		}

		manifest.WorkingDir = pluginPath
		plugins = append(plugins, &SubprocessCleaner{manifest: manifest})
	}
	return plugins
}