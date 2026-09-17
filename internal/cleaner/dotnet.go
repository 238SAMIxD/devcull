package cleaner

import (
	"os/exec"
	"strings"
)

type DotnetCleaner struct{}

func (d *DotnetCleaner) Name() string {
	return "Dotnet"
}

func (d *DotnetCleaner) Category() string {
	return "Languages"
}

func (d *DotnetCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return false
	}
	if err := exec.Command("dotnet", "--version").Run(); err != nil {
		return false
	}
	return true
}

func (d *DotnetCleaner) getCachePaths() []string {
	out, err := exec.Command("dotnet", "nuget", "locals", "all", "--list").Output()
	if err != nil {
		return nil
	}

	var paths []string
	lines := strings.Split(string(out), "\n")
	
	prefixes := []string{"http-cache:", "global-packages:", "temp:", "plugins-cache:"}

	for _, line := range lines {
		for _, prefix := range prefixes {
			if idx := strings.Index(line, prefix); idx != -1 {
				path := strings.TrimSpace(line[idx+len(prefix):])
				if path != "" {
					paths = append(paths, path)
				}
				break
			}
		}
	}
	return paths
}

func (d *DotnetCleaner) EstimateReclaimable() (int64, error) {
	if !d.IsInstalled() {
		return 0, ErrToolNotInstalled
	}

	paths := d.getCachePaths()
	var total int64
	for _, p := range paths {
		size, _ := dirSize(p)
		total += size
	}

	return total, nil
}

func (d *DotnetCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := d.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun || reclaimable == 0 {
		return reclaimable, nil
	}

	if err := exec.Command("dotnet", "nuget", "locals", "all", "--clear").Run(); err != nil {
		return 0, err
	}

	return reclaimable, nil
}