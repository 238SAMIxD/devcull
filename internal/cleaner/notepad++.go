package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)
type NotepadPlusPlusCleaner struct{}

func (c *NotepadPlusPlusCleaner) Name() string       { return "Notepad++" }
func (c *NotepadPlusPlusCleaner) Category() Category { return CategoryIDE }

func (c *NotepadPlusPlusCleaner) getPaths() []string {
	if runtime.GOOS == "windows" {
		return []string{filepath.Join(os.Getenv("APPDATA"), "Notepad++", "backup")}
	}
	return nil
}

func (c *NotepadPlusPlusCleaner) IsInstalled() bool {
	for _, p := range c.getPaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}
func (c *NotepadPlusPlusCleaner) EstimateReclaimable() (int64, error) {
	return dirsSize(c.getPaths()), nil
}

func (c *NotepadPlusPlusCleaner) Clean(dryRun bool) (int64, error) {
	return cleanDirs(c.getPaths(), dryRun)
}
