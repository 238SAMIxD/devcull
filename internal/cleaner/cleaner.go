package cleaner

import (
	"os"
	"path/filepath"
	"strings"
)

// Cleaner defines a tool whose caches can be scanned and cleaned.
// Callers must check IsInstalled() before calling EstimateReclaimable() or Clean().
type Cleaner interface {
	Name() string
	Category() string
	IsInstalled() bool
	EstimateReclaimable() (int64, error)
	Clean(dryRun bool) (int64, error)
}

func ResolveAlias(name string) string {
	switch strings.ToLower(name) {
	case "brew":
		return "homebrew"
	case "python", "python3", "pip3":
		return "pip"
	case "node", "nodejs":
		return "npm"
	case "rust":
		return "cargo"
	case "golang":
		return "go"
	case "mac", "apple", "ios", "macos":
		return "cocoapods"
	case "cs", "c#", "csharp", ".net", "nuget":
		return "dotnet"
	default:
		return strings.ToLower(name)
	}
}

func All() []Cleaner {
	return []Cleaner{
			&NvmCleaner{},
			
			&NpmCleaner{},
			&PnpmCleaner{},
			&YarnCleaner{},

		 	&GoCleaner{},
			&CargoCleaner{},

			&BrewCleaner{},
			
			&PipCleaner{},
			&UvCleaner{},
			&PoetryCleaner{},

			&DockerCleaner{},

			&BunCleaner{},
			&DenoCleaner{},

			&DotnetCleaner{},
		}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size, err
}