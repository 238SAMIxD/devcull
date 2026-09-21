package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"syscall"
	"time"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

type Manifest struct {
	Name     string           `json:"name"`
	Category cleaner.Category `json:"category"`
	Aliases  []string         `json:"aliases"`

	Description string   `json:"description"`
	Version     string   `json:"version"`
	Entrypoint  []string `json:"entrypoint"`
	WorkingDir  string   `json:"-"`
}

type SubprocessCleaner struct {
	manifest Manifest
}

func (p *SubprocessCleaner) Name() string               { return p.manifest.Name }
func (p *SubprocessCleaner) Category() cleaner.Category { return p.manifest.Category }
func (p *SubprocessCleaner) Aliases() []string          { return p.manifest.Aliases }

func (p *SubprocessCleaner) IsInstalled(ctx context.Context) bool {
	if len(p.manifest.Entrypoint) == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmdArgs := append([]string{}, p.manifest.Entrypoint[1:]...)
	cmdArgs = append(cmdArgs, "installed")
	cmd := exec.CommandContext(ctx, p.manifest.Entrypoint[0], cmdArgs...)
	cmd.Dir = p.manifest.WorkingDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

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

func (p *SubprocessCleaner) EstimateReclaimable(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmdArgs := append([]string{}, p.manifest.Entrypoint[1:]...)
	cmdArgs = append(cmdArgs, "estimate")
	cmd := exec.CommandContext(ctx, p.manifest.Entrypoint[0], cmdArgs...)
	cmd.Dir = p.manifest.WorkingDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

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

func (p *SubprocessCleaner) Clean(ctx context.Context, dryRun bool) (int64, error) {
	args := append([]string{}, p.manifest.Entrypoint[1:]...)
	args = append(args, "clean")
	if dryRun {
		args = append(args, "--dry-run")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, p.manifest.Entrypoint[0], args...)
	cmd.Dir = p.manifest.WorkingDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

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
