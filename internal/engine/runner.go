package engine

import (
	"context"
	"runtime"
	"sync"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

type Result struct {
	CleanerName string
	Category    cleaner.Category
	Reclaimed   int64
	Skipped     bool
	Err         error
}

type ScanResult struct {
	CleanerName string
	Category    cleaner.Category
	Reclaimable int64
	Skipped     bool
	Err         error
}

func execute[T any](
	ctx context.Context,
	cleaners []cleaner.Cleaner,
	concurrencyLimit int,
	makeCancelResult func(c cleaner.Cleaner, err error) T,
	makeSkippedResult func(c cleaner.Cleaner) T,
	doWork func(ctx context.Context, c cleaner.Cleaner) T,
	onProgress func(),
) []T {
	var results []T
	var wg sync.WaitGroup
	var mu sync.Mutex

	var sem chan struct{}
	if concurrencyLimit > 0 {
		sem = make(chan struct{}, concurrencyLimit)
	}

	for _, c := range cleaners {
		wg.Add(1)
		go func(clr cleaner.Cleaner) {
			defer wg.Done()
			if onProgress != nil {
				defer onProgress()
			}

			if sem != nil {
				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					mu.Lock()
					results = append(results, makeCancelResult(clr, ctx.Err()))
					mu.Unlock()
					return
				}
				defer func() { <-sem }()
			}

			select {
			case <-ctx.Done():
				mu.Lock()
				results = append(results, makeCancelResult(clr, ctx.Err()))
				mu.Unlock()
				return
			default:
			}

			if !clr.IsInstalled(ctx) {
				mu.Lock()
				results = append(results, makeSkippedResult(clr))
				mu.Unlock()
				return
			}

			res := doWork(ctx, clr)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	return results
}

func Run(ctx context.Context, cleaners []cleaner.Cleaner, dryRun bool, onProgress func()) []Result {
	return execute(ctx, cleaners, runtime.NumCPU(),
		func(c cleaner.Cleaner, err error) Result {
			return Result{CleanerName: c.Name(), Category: c.Category(), Err: err}
		},
		func(c cleaner.Cleaner) Result {
			return Result{CleanerName: c.Name(), Category: c.Category(), Skipped: true}
		},
		func(ctx context.Context, c cleaner.Cleaner) Result {
			reclaimed, err := c.Clean(ctx, dryRun)
			return Result{CleanerName: c.Name(), Category: c.Category(), Reclaimed: reclaimed, Err: err}
		},
		onProgress,
	)
}

func Scan(ctx context.Context, cleaners []cleaner.Cleaner, onProgress func()) []ScanResult {
	return execute(ctx, cleaners, runtime.NumCPU(),
		func(c cleaner.Cleaner, err error) ScanResult {
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Err: err}
		},
		func(c cleaner.Cleaner) ScanResult {
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Skipped: true}
		},
		func(ctx context.Context, c cleaner.Cleaner) ScanResult {
			reclaimable, err := c.EstimateReclaimable(ctx)
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Reclaimable: reclaimable, Err: err}
		},
		onProgress,
	)
}

func RunPlugins(ctx context.Context, cleaners []cleaner.Cleaner, dryRun bool, onProgress func()) []Result {
	return execute(ctx, cleaners, runtime.NumCPU()*4,
		func(c cleaner.Cleaner, err error) Result {
			return Result{CleanerName: c.Name(), Category: c.Category(), Err: err}
		},
		func(c cleaner.Cleaner) Result {
			return Result{CleanerName: c.Name(), Category: c.Category(), Skipped: true}
		},
		func(ctx context.Context, c cleaner.Cleaner) Result {
			reclaimed, err := c.Clean(ctx, dryRun)
			return Result{CleanerName: c.Name(), Category: c.Category(), Reclaimed: reclaimed, Err: err}
		},
		onProgress,
	)
}

func ScanPlugins(ctx context.Context, cleaners []cleaner.Cleaner, onProgress func()) []ScanResult {
	return execute(ctx, cleaners, runtime.NumCPU()*4,
		func(c cleaner.Cleaner, err error) ScanResult {
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Err: err}
		},
		func(c cleaner.Cleaner) ScanResult {
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Skipped: true}
		},
		func(ctx context.Context, c cleaner.Cleaner) ScanResult {
			reclaimable, err := c.EstimateReclaimable(ctx)
			return ScanResult{CleanerName: c.Name(), Category: c.Category(), Reclaimable: reclaimable, Err: err}
		},
		onProgress,
	)
}
