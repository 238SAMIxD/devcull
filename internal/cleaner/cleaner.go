package cleaner

import (
	"os"
	"path/filepath"
	"strings"
)



type Category string

const (
	CategoryNode   Category = "Node"
	CategoryPython Category = "Python"
	CategoryJava   Category = "Java"
	CategoryDotnet Category = ".NET"
	CategoryGo     Category = "Go"
	CategoryRust   Category = "Rust"
	CategoryApple  Category = "Apple Ecosystem"
	CategorySystem Category = "System & DevOps"
)

func AllCategories() []Category {
	return []Category{
		CategoryNode,
		CategoryPython,
		CategoryJava,
		CategoryDotnet,
		CategoryGo,
		CategoryRust,
		CategoryApple,
		CategorySystem,
	}
}

type Cleaner interface {
	Name() string
	Category() Category
	IsInstalled() bool
	EstimateReclaimable() (int64, error)
	Clean(dryRun bool) (int64, error)
}

func All() []Cleaner {
	return []Cleaner{
			// Node.js
			&NvmCleaner{},
			&NpmCleaner{},
			&PnpmCleaner{},
			&YarnCleaner{},
			&BunCleaner{},
			&DenoCleaner{},

		 	// Go
			&GoCleaner{},

			// Rust
			&CargoCleaner{},

			// System
			&BrewCleaner{},
			&DockerCleaner{},
			
			// Python
			&PipCleaner{},
			&UvCleaner{},
			&PoetryCleaner{},

			// .NET
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

func MatchesArg(cleaner Cleaner, arg string) bool {
	argLower := strings.ToLower(strings.TrimSpace(arg))
	nameLower := strings.ToLower(cleaner.Name())
	catLower := strings.ToLower(string(cleaner.Category()))

	if nameLower == argLower {
		return true
	}

	if catLower == argLower || strings.Contains(catLower, argLower) {
		return true
	}

	switch argLower {
	case "brew", "macos":
		return nameLower == "homebrew"
	case "py", "python3":
		return cleaner.Category() == CategoryPython
	case "node", "nodejs", "js", "ts", "javascript":
		return cleaner.Category() == CategoryNode
	case "java", "jvm":
		return cleaner.Category() == CategoryJava
	case "cs", "c#", "csharp", ".net", "nuget":
		return cleaner.Category() == CategoryDotnet
	case "golang":
		return cleaner.Category() == CategoryGo
	case "mac", "apple", "ios":
		return nameLower == "cocoapods" || cleaner.Category() == CategoryApple
	}

	return false
}

