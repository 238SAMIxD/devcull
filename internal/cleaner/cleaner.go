package cleaner

import (
	"errors"
	"strings"
)


var ErrToolNotInstalled = errors.New("tool not installed")

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