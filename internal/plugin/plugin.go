package plugin

import (
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

type Manifest struct {
	Name        string           `json:"name"`
	Category    cleaner.Category `json:"category"`
	Description string           `json:"description"`
	Version     string           `json:"version"`
	Entrypoint  []string         `json:"entrypoint"`
	WorkingDir  string           `json:"-"`
}

type SubprocessCleaner struct {
	manifest Manifest
}

func (p *SubprocessCleaner) Name() string               { return p.manifest.Name }
func (p *SubprocessCleaner) Category() cleaner.Category { return p.manifest.Category }

func (p *SubprocessCleaner) IsInstalled() bool {
	if len(p.manifest.Entrypoint) == 0 {
		return false
	}

	cmdArgs := append(p.manifest.Entrypoint[1:], "installed")
	cmd := exec.Command(p.manifest.Entrypoint[0], cmdArgs...)
	cmd.Dir = p.manifest.WorkingDir

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	var res struct {
		Installed bool `json:"installed"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		return false
	}
	return res.Installed
}

func (p *SubprocessCleaner) EstimateReclaimable() (int64, error) {
	cmdArgs := append(p.manifest.Entrypoint[1:], "estimate")
	cmd := exec.Command(p.manifest.Entrypoint[0], cmdArgs...)
	cmd.Dir = p.manifest.WorkingDir

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var res struct {
		ReclaimableBytes int64  `json:"reclaimable_bytes"`
		Error            string `json:"error"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		return 0, err
	}

	if res.Error != "" {
		return 0, errors.New(res.Error)
	}
	return res.ReclaimableBytes, nil
}

func (p *SubprocessCleaner) Clean(dryRun bool) (int64, error) {
	args := append(p.manifest.Entrypoint[1:], "clean")
	if dryRun {
		args = append(args, "--dry-run")
	}

	cmd := exec.Command(p.manifest.Entrypoint[0], args...)
	cmd.Dir = p.manifest.WorkingDir

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var res struct {
		ReclaimedBytes int64  `json:"reclaimed_bytes"`
		Error          string `json:"error"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		return 0, err
	}

	if res.Error != "" {
		return 0, errors.New(res.Error)
	}
	return res.ReclaimedBytes, nil
}
