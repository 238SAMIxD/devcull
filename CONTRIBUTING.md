# Contributing to devcull

We welcome contributions! Whether you are adding a native Go cleaner to the core engine or building a community plugin, here is how you can help.

## 1. Adding a Native Go Cleaner

If you want to add a highly popular, standard developer tool to the core engine, you can write a native Go cleaner.

1. Create a new file in `internal/cleaner/` (e.g., `rust.go`).
2. Implement the `Cleaner` interface:

   ```go
   type Cleaner interface {
       Name() string
       Category() Category
       IsInstalled() bool
       EstimateReclaimable() (int64, error)
       Clean(dryRun bool) (int64, error)
   }
   ```

3. Register your struct in the `Native()` array inside `internal/cleaner/cleaner.go`.
4. Open a Pull Request!

## 2. Building a Custom Plugin

`devcull` supports language-agnostic plugins. The engine will execute your script as a subprocess and read the JSON output.

To create a plugin, make a folder in `~/.config/devcull/plugins/` (or `./plugins` locally) containing two things:

**1. `manifest.json**`

```json
{
  "name": "My Custom Cleaner",
  "category": "Java",
  "entrypoint": ["bash", "cleaner.sh"]
}
```

**2. Your Executable (`cleaner.sh` example)**
Your script just needs to read the command-line argument (`installed`, `estimate`, or `clean`) and print valid JSON to `stdout`.

```bash
#!/bin/bash
COMMAND=$1

case "$COMMAND" in
  "installed")
    echo '{"installed": true}'
    ;;
  "estimate")
    # Logic to calculate bytes...
    echo '{"reclaimable_bytes": 1048576, "error": ""}'
    ;;
  "clean")
    # Logic to delete files...
    echo '{"reclaimed_bytes": 1048576, "error": ""}'
    ;;
esac
```

## Running Tests

Before submitting a PR, ensure all CI checks will pass by running:

```bash
go test -v -race ./...
go build -v ./cmd/devcull
```
