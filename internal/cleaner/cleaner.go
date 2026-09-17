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
	CategoryCpp    Category = "C++"
	CategoryCSharp Category = "C#"
	CategoryGo     Category = "Go"
	CategoryRust   Category = "Rust"
	CategoryFlutter	 Category = "Flutter"
	CategoryApple  Category = "Apple Ecosystem"
	CategorySystem Category = "System & DevOps"
)

func AllCategories() []Category {
	return []Category{
		CategoryNode,
		CategoryPython,
		CategoryJava,
		CategoryCpp,
		CategoryCSharp,
		CategoryGo,
		CategoryRust,
		CategoryFlutter,
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

			// C/C++
			&UnrealCleaner{},
			&CcacheCleaner{},
			&ConanCleaner{},

			// C#
			&DotnetCleaner{},
			&UnityCleaner{},

			// Java
      &MavenCleaner{},
      &GradleCleaner{},
      &KotlinCleaner{},
      &AndroidCleaner{},

			// Flutter
			&DartCleaner{},
			&FvmCleaner{},

			// Apple
			&CocoaPodsCleaner{},
			&SwiftPMCleaner{},
			&XcodeCleaner{},
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
	case "py", "python3", "python":
		return cleaner.Category() == CategoryPython
	case "node", "node.js", "nodejs", "javascript", "typescript", "js", "ts":
		return cleaner.Category() == CategoryNode
	case "java", "jvm":
		return cleaner.Category() == CategoryJava
	case "cs", "c#", "csharp", ".net":
		return cleaner.Category() == CategoryCSharp
	case "golang":
		return cleaner.Category() == CategoryGo
	case "mac", "apple", "ios", "macos", "darwin", "pods":
		return nameLower == "cocoapods" || cleaner.Category() == CategoryApple
	case "c", "cpp", "c++":
		return cleaner.Category() == CategoryCpp
	case "flutter", "pub":
		return cleaner.Category() == CategoryFlutter
		
	case "brew":
		return nameLower == "homebrew"
	case "pip3":
		return nameLower == "pip"
	case "nuget":
		return nameLower == "dotnet"
	case "mvn":
		return nameLower == "maven"
	case "kt", "konan":
		return nameLower == "kotlin"
	case "ue", "ue5", "unreal", "unrealengine":
		return nameLower == "unreal engine"
	case "swift", "spm":
		return nameLower == "swiftpm"
	case "x-code", "deriveddata", "commandlinetools", "clt", "cmdlinetools":
		return nameLower == "xcode"
	}

	return false
}

