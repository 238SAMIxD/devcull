# devcull 🧹

> A lightning-fast, extensible CLI to reclaim gigabytes of stale developer caches.

`devcull` safely hunts down and removes massive, forgotten cache directories left behind by IDEs, package managers, and toolchains. It features a concurrent Go execution engine and a language-agnostic plugin system, allowing anyone to extend it.

## Features

- **Blazing Fast:** Concurrent engine scans your entire filesystem in milliseconds.
- **Safe by Default:** Strictly targets disposable caches (indexes, derived data, cached packages). Never touches user settings or source code.
- **Extensible Architecture:** Write custom cleaners in _any_ language (TypeScript, Bash, Python, Go) using standard subprocesses and JSON.
- **Dry-Run Mode:** See exactly what will be deleted and how much space you will save before touching the disk.

## Installation

**Requires Go 1.27.1+**

```bash
go install github.com/238SAMIxD/devcull/cmd/devcull@latest
```

_(Note: Pre-compiled binaries for Mac, Linux, and Windows are coming in the official v0.1.0 release!)_

## Usage

**1. Scan for reclaimable space**
Find out exactly how much disk space is being wasted by stale caches.

```bash
devcull scan
```

**2. Clean caches**
Reclaim the space. (You can also run `devcull clean --dry-run` to test the waters).

```bash
devcull clean
```

## Supported Targets (Native)

- **Apple Ecosystem:** Xcode (DerivedData), SwiftPM, iOS Simulator Caches
- **Java:** Android Studio, Maven, Gradle, Kotlin
- **Node & JS:** npm, pnpm, yarn, bun, deno caches
- **IDEs & Editors:** VS Code, JetBrains (IntelliJ, WebStorm, etc.), Neovim, Visual Studio, Eclipse, Sublime Text

## The Plugin Ecosystem

`devcull` doesn't just rely on native Go implementations. You can drop a custom plugin into `~/.config/devcull/plugins/` written in **any language**. The engine communicates with plugins by passing commands as positional CLI arguments (e.g., `./plugin clean`) and expects a single JSON response on standard output.

Want to build a plugin to clear out Chrome's cache in TypeScript, or a Docker cleanup script in Bash? Check out [Contributing Guide](CONTRIBUTING.md).
