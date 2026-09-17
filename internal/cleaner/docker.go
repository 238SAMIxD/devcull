package cleaner

import (
	"os/exec"
	"strings"

	"github.com/238SAMIxD/devcull/internal/utils"
)

type DockerCleaner struct{}

func (d *DockerCleaner) Name() string {
	return "Docker"
}

func (d *DockerCleaner) Category() string {
	return "Containers"
}

func (d *DockerCleaner) IsInstalled() bool {
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	if err := exec.Command("docker", "info").Run(); err != nil {
		return false
	}
	return true
}

func (d *DockerCleaner) EstimateReclaimable() (int64, error) {
	out, err := exec.Command("docker", "system", "df", "--format", "{{.Reclaimable}}").Output()
	if err != nil {
		return 0, err
	}

	var totalEstimate int64
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		totalEstimate += utils.ParseByteString(line)
	}

	return totalEstimate, nil
}

func (d *DockerCleaner) Clean(dryRun bool) (int64, error) {
	reclaimable, err := d.EstimateReclaimable()
	if err != nil {
		return 0, err
	}

	if dryRun {
		return reclaimable, nil
	}

	out, err := exec.Command("docker", "system", "prune", "-f").Output()
	if err != nil {
		return 0, err
	}

	return utils.ParseByteString(string(out)), nil
}