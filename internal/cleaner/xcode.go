package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
)

type XcodeCleaner struct{}

func (x *XcodeCleaner) Name() string { return "Xcode" }
func (x *XcodeCleaner) Category() Category { return CategoryApple }

func (x *XcodeCleaner) getCachePaths() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	
	return []string{
		filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData"),
		filepath.Join(home, "Library", "Developer", "Xcode", "iOS DeviceSupport"),
		filepath.Join(home, "Library", "Developer", "Xcode", "watchOS DeviceSupport"),
		filepath.Join(home, "Library", "Developer", "Xcode", "tvOS DeviceSupport"),
	}
}

func (x *XcodeCleaner) IsInstalled() bool {
	for _, p := range x.getCachePaths() {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

func (x *XcodeCleaner) EstimateReclaimable() (int64, error) {
	var total int64
	for _, p := range x.getCachePaths() {
		size, _ := dirSize(p)
		total += size
	}
	return total, nil
}

func (x *XcodeCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := x.EstimateReclaimable()
	if err != nil || dryRun || reclaimable == 0 {
		return reclaimable, err
	}
	for _, p := range x.getCachePaths() {
		_ = os.RemoveAll(p)
	}
	return reclaimable, nil
}