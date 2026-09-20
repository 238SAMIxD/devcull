package cleaner

import (
	"context"
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
	Aliases() []string

	IsInstalled(ctx context.Context) bool
	EstimateReclaimable(ctx context.Context) (int64, error)
	Clean(ctx context.Context, dryRun bool) (int64, error)
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

func dirSize(ctx context.Context, path string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("empty path provided to dirSize")
	}
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
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

func dirsSize(ctx context.Context, paths []string) (int64, error) {
	var total int64
	var firstErr error
	for _, p := range paths {
		size, err := dirSize(ctx, p)
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

	var exactMatchRoots []string
	var prefixMatchTrees []string

	home, err := os.UserHomeDir()
	if err == nil {
		cleanHome := filepath.Clean(home)
		exactMatchRoots = append(exactMatchRoots, cleanHome)
		prefixMatchTrees = append(prefixMatchTrees,
			filepath.Join(cleanHome, "Desktop"),
			filepath.Join(cleanHome, "Documents"),
			filepath.Join(cleanHome, "Downloads"),
			filepath.Join(cleanHome, "Pictures"),
			filepath.Join(cleanHome, "Music"),
			filepath.Join(cleanHome, "Movies"),
			filepath.Join(cleanHome, "Public"),
		)
	}

	if os.PathSeparator == '/' {
		exactMatchRoots = append(exactMatchRoots, "/Users")
		prefixMatchTrees = append(prefixMatchTrees,
			"/usr", "/usr/bin", "/usr/lib", "/usr/local", "/usr/local/bin",
			"/bin", "/sbin", "/etc", "/var", "/tmp", "/Library", "/System",
			"/Applications", "/Network", "/Volumes",
		)
		if err == nil {
			prefixMatchTrees = append(prefixMatchTrees, filepath.Join(home, "Library"))
		}
	} else if os.PathSeparator == '\\' {
		sysRoot := os.Getenv("SystemRoot")
		if sysRoot == "" {
			sysRoot = `C:\Windows`
		}
		progFiles := os.Getenv("ProgramFiles")
		if progFiles == "" {
			progFiles = `C:\Program Files`
		}
		progFiles86 := os.Getenv("ProgramFiles(x86)")
		if progFiles86 == "" {
			progFiles86 = `C:\Program Files (x86)`
		}
		exactMatchRoots = append(exactMatchRoots, `C:\Users`)
		prefixMatchTrees = append(prefixMatchTrees,
			sysRoot,
			filepath.Join(sysRoot, "System32"),
			progFiles,
			progFiles86,
		)
		if err == nil {
			prefixMatchTrees = append(prefixMatchTrees, filepath.Join(home, "AppData"))
		}
	}

	for _, root := range exactMatchRoots {
		if root == "" || root == "." {
			continue
		}
		cleanRoot := filepath.Clean(root)
		if strings.EqualFold(abs, cleanRoot) {
			return false
		}
	}

	for _, tree := range prefixMatchTrees {
		if tree == "" || tree == "." {
			continue
		}
		cleanTree := filepath.Clean(tree)
		if strings.EqualFold(abs, cleanTree) {
			return false
		}
		if strings.HasPrefix(strings.ToLower(abs), strings.ToLower(cleanTree+string(filepath.Separator))) {
			return false
		}
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

func cleanDirs(ctx context.Context, paths []string, dryRun bool) (int64, error) {
	before, err := dirsSize(ctx, paths)
	if err != nil {
		return 0, err
	}
	if before == 0 || dryRun {
		return before, nil
	}

	var firstErr error
	for _, p := range paths {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		if !isSafeToDelete(p) {
			continue
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		abs = filepath.Clean(abs)
		if err := os.RemoveAll(abs); err != nil && firstErr == nil {
			if !os.IsNotExist(err) {
				firstErr = err
			}
		}
	}

	after, err := dirsSize(ctx, paths)
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
	argLower := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(arg, " ", ""), "-", ""))
	nameLower := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(cleaner.Name(), " ", ""), "-", ""))
	catLower := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(string(cleaner.Category()), " ", ""), "-", ""))

	if nameLower == argLower {
		return true
	}

	if catLower == argLower {
		return true
	}

	for _, alias := range cleaner.Aliases() {
		aliasLower := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(alias, " ", ""), "-", ""))
		if aliasLower == argLower {
			return true
		}
	}

	switch argLower {
	case "py", "python3", "python":
		return cleaner.Category() == CategoryPython
	case "node", "nodejs", "javascript", "typescript", "js", "ts":
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
		return nameLower == "unrealengine"
	case "swift", "spm":
		return nameLower == "swiftpm"
	case "xcode", "deriveddata", "commandlinetools", "clt", "cmdlinetools":
		return nameLower == "xcode"
	}

	return false
}
