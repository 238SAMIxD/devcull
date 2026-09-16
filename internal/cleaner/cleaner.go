package cleaner

import "errors"

var ErrToolNotInstalled = errors.New("tool not installed")

type Cleaner interface {
	Name() string
	Category() string
	IsInstalled() bool
	EstimateReclaimable() (int64, error)
	Clean(dryRun bool) (int64, error)
}