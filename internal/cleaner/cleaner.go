package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Category string

const (
	CategoryNode    Category = "Node"
	CategoryPython  Category = "Python"
	CategoryJava    Category = "Java"
	CategoryCpp     Category = "C++"
	CategoryCSharp  Category = "C#"
	CategoryGo      Category = "Go"
	CategoryRust    Category = "Rust"
	CategoryFlutter Category = "Flutter"
	CategoryPHP     Category = "PHP"
	CategoryApple   Category = "Apple Ecosystem"
	CategoryIDE     Category = "IDE"
	CategorySystem  Category = "System & DevOps"
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
		CategoryPHP,
		CategoryApple,
		CategoryIDE,
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

func Native() []Cleaner {
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

		// PHP
		&ComposerCleaner{},
		&PhpbrewCleaner{},

		// IDE
		&EclipseCleaner{},
		&JetBrainsCleaner{},
		&VSCodeCleaner{},
		&NeovimCleaner{},
		&NetBeansCleaner{},
		&SublimeCleaner{},
		&VisualStudioCleaner{},
	}
}

func dirSize(path string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("empty path provided to dirSize")
	}
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			size += info.Size()
		}
		return nil
	})
	if os.IsNotExist(err) {
		return size, nil
	}
	return size, err
}

func dirsSize(paths []string) (int64, error) {
	var total int64
	var firstErr error
	for _, p := range paths {
		size, err := dirSize(p)
		if err != nil && firstErr == nil {
			firstErr = err
		}
		total += size
	}
	return total, firstErr
}

func isSafeToDelete(targetPath string) bool {
	if targetPath == "" {
		return false
	}

	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return false
	}
	abs = filepath.Clean(abs)

	vol := filepath.VolumeName(abs)
	if abs == "/" || abs == "\\" || abs == vol+"\\" || abs == vol+"/" {
		return false
	}

	home, err := os.UserHomeDir()
	if err == nil && abs == filepath.Clean(home) {
		return false
	}

	pWithoutVol := abs[len(vol):]
	parts := strings.Split(filepath.ToSlash(pWithoutVol), "/")

	depth := 0
	for _, part := range parts {
		if part != "" {
			depth++
		}
	}

	if depth < 2 {
		return false
	}

	return true
}

func cleanDirs(paths []string, dryRun bool) (int64, error) {
	before, err := dirsSize(paths)
	if err != nil {
		return 0, err
	}
	if before == 0 || dryRun {
		return before, nil
	}

	var firstErr error
	for _, p := range paths {
		if !isSafeToDelete(p) {
			continue
		}
		if err := os.RemoveAll(p); err != nil && firstErr == nil {
			if !os.IsNotExist(err) {
				firstErr = err
			}
		}
	}

	after, err := dirsSize(paths)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
		return 0, firstErr
	}

	reclaimed := before - after
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, firstErr
}

func MatchesArg(cleaner Cleaner, arg string) bool {
	argLower := strings.ToLower(strings.TrimSpace(arg))
	nameLower := strings.ToLower(cleaner.Name())
	catLower := strings.ToLower(string(cleaner.Category()))

	if nameLower == argLower {
		return true
	}

	if catLower == argLower {
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
	case "mac", "apple", "ios", "macos", "darwin":
		return cleaner.Category() == CategoryApple
	case "c", "cpp", "c++":
		return cleaner.Category() == CategoryCpp
	case "flutter":
		return cleaner.Category() == CategoryFlutter
	case "php":
		return cleaner.Category() == CategoryPHP
	case "editor", "editors", "ides", "codeeditor", "texteditor", "text", "environment", "env", "dev":
		return cleaner.Category() == CategoryIDE
	case "system", "devops", "ops":
		return cleaner.Category() == CategorySystem

	case "pods":
		return nameLower == "cocoapods"
	case "pub":
		return nameLower == "dart"
	case "composer", "laravel":
		return nameLower == "composer"
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
