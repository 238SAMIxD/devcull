package cleaner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
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
		&ArduinoCleaner{},
		&ArduinoLabCleaner{},
	}
}

var (
	sem       = make(chan struct{}, runtime.NumCPU()*8)
	workerSem = make(chan struct{}, runtime.NumCPU()*64)
)

func walkDirConcurrent(ctx context.Context, path string, size *atomic.Int64, wg *sync.WaitGroup, errMu *sync.Mutex, firstErr *error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	sem <- struct{}{}
	entries, err := os.ReadDir(path)
	<-sem

	if err != nil {
		if !os.IsNotExist(err) {
			errMu.Lock()
			if *firstErr == nil {
				*firstErr = err
			}
			errMu.Unlock()
		}
		return
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			select {
			case workerSem <- struct{}{}:
				wg.Add(1)
				go func() {
					defer func() { <-workerSem }()
					walkDirConcurrent(ctx, subPath, size, wg, errMu, firstErr)
				}()
			default:
				wg.Add(1)
				walkDirConcurrent(ctx, subPath, size, wg, errMu, firstErr)
			}
		} else {
			info, err := entry.Info()
			if err == nil {
				size.Add(info.Size())
			} else if !os.IsNotExist(err) {
				errMu.Lock()
				if *firstErr == nil {
					*firstErr = err
				}
				errMu.Unlock()
			}
		}
	}
}

func dirSize(ctx context.Context, path string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("empty path provided to dirSize")
	}

	info, err := os.Lstat(path)
	if err != nil {
		if errCtx := ctx.Err(); errCtx != nil {
			return 0, errCtx
		}
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	var size atomic.Int64
	if !info.IsDir() {
		if errCtx := ctx.Err(); errCtx != nil {
			return 0, errCtx
		}
		size.Add(info.Size())
		return size.Load(), nil
	}

	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	wg.Add(1)
	go walkDirConcurrent(ctx, path, &size, &wg, &errMu, &firstErr)
	wg.Wait()

	if ctx.Err() != nil {
		return size.Load(), ctx.Err()
	}

	return size.Load(), firstErr
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

	evalPath, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			evalPath = targetPath
		} else {
			return false
		}
	}

	abs, err := filepath.Abs(evalPath)
	if err != nil {
		return false
	}
	abs = filepath.Clean(abs)

	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	cleanHome := filepath.Clean(home)

	var safeRoots []string

	if cacheDir, err := os.UserCacheDir(); err == nil {
		safeRoots = append(safeRoots, filepath.Clean(cacheDir))
	}

	safeRoots = append(safeRoots,
		filepath.Join(cleanHome, "Library", "Caches"),
		filepath.Join(cleanHome, "Library", "Developer"),
		filepath.Join(cleanHome, ".cache"),
		filepath.Join(cleanHome, ".local", "share"),
		filepath.Join(cleanHome, ".gradle"),
		filepath.Join(cleanHome, ".bun"),
		filepath.Join(cleanHome, ".nvm"),
		filepath.Join(cleanHome, ".cargo"),
		filepath.Join(cleanHome, ".rustup"),
		filepath.Join(cleanHome, ".npm"),
		filepath.Join(cleanHome, ".pnpm-state"),
		filepath.Join(cleanHome, ".yarn"),
		filepath.Join(cleanHome, ".docker"),
		filepath.Join(cleanHome, ".poetry"),
		filepath.Join(cleanHome, ".ccache"),
		filepath.Join(cleanHome, ".conan"),
		filepath.Join(cleanHome, ".conan2"),
		filepath.Join(cleanHome, ".dotnet"),
		filepath.Join(cleanHome, ".android"),
		filepath.Join(cleanHome, ".dartServer"),
		filepath.Join(cleanHome, ".fvm"),
		filepath.Join(cleanHome, ".cocoapods"),
		filepath.Join(cleanHome, ".phpbrew"),
		filepath.Join(cleanHome, ".eclipse"),
		filepath.Join(cleanHome, ".vscode"),
	)

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		safeRoots = append(safeRoots, filepath.Clean(localAppData))
	}

	if appData := os.Getenv("APPDATA"); appData != "" {
		safeRoots = append(safeRoots, filepath.Clean(appData))
	}

	if tempDir := os.Getenv("TEMP"); tempDir != "" {
		safeRoots = append(safeRoots, filepath.Clean(tempDir))
	}

	if configDir, err := os.UserConfigDir(); err == nil {
		safeRoots = append(safeRoots, filepath.Join(configDir, "arduino-ide"))
		safeRoots = append(safeRoots, filepath.Join(configDir, "arduino-lab-for-micropython"))
	}
	safeRoots = append(safeRoots, filepath.Join(cleanHome, "Library", "Application Support", "arduino-ide"))
	safeRoots = append(safeRoots, filepath.Join(cleanHome, "Library", "Application Support", "arduino-lab-for-micropython"))

	for _, root := range safeRoots {
		if root == "" || root == "." {
			continue
		}
		cleanRoot := filepath.Clean(root)
		if strings.EqualFold(abs, cleanRoot) {
			continue
		}
		if strings.HasPrefix(strings.ToLower(abs), strings.ToLower(cleanRoot+string(filepath.Separator))) {
			return true
		}
	}

	return false
}

func removeAllCtx(ctx context.Context, path string) error {
	var dirs []string
	err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
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
		if d.IsDir() {
			dirs = append(dirs, p)
			return nil
		}
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := os.Remove(dirs[i]); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func cleanDirs(ctx context.Context, paths []string, dryRun bool) (int64, error) {
	var safePaths []string
	for _, p := range paths {
		if isSafeToDelete(p) {
			safePaths = append(safePaths, p)
		} else {
			fmt.Fprintf(os.Stderr, "WARN: Path rejected by safety guards: %s\n", p)
		}
	}

	before, err := dirsSize(ctx, safePaths)
	if err != nil {
		return 0, err
	}
	if before == 0 || dryRun {
		return before, nil
	}

	var firstErr error
	for _, p := range safePaths {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		abs = filepath.Clean(abs)
		if err := removeAllCtx(ctx, abs); err != nil && firstErr == nil {
			if !os.IsNotExist(err) {
				firstErr = err
			}
		}
	}

	after, err := dirsSize(ctx, safePaths)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
		after = 0
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
