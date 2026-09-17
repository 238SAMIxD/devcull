package engine

import (
	"sync"

	"github.com/238SAMIxD/devcull/internal/cleaner"
)

type Result struct {
	CleanerName string
	Reclaimed   int64
	Err         error
}

type ScanResult struct {
	CleanerName string
	Category    cleaner.Category
	Reclaimable int64
	Err         error
}

func Run(cleaners []cleaner.Cleaner, dryRun bool) []Result {
	var wg sync.WaitGroup
	resultsCh := make(chan Result, len(cleaners))

	for _, c := range cleaners {
		wg.Add(1)
		go func(clr cleaner.Cleaner) {
			defer wg.Done()

			if !clr.IsInstalled() {
				return
			}

			reclaimed, err := clr.Clean(dryRun)
			resultsCh <- Result{
				CleanerName: clr.Name(),
				Reclaimed:   reclaimed,
				Err:         err,
			}
		}(c)
	}

	wg.Wait()
	close(resultsCh)

	var results []Result
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}

func Scan(cleaners []cleaner.Cleaner) []ScanResult {
	var wg sync.WaitGroup
	resultsCh := make(chan ScanResult, len(cleaners))

	for _, c := range cleaners {
		wg.Add(1)
		go func(clr cleaner.Cleaner) {
			defer wg.Done()

			if !clr.IsInstalled() {
				return
			}

			reclaimable, err := clr.EstimateReclaimable()
			resultsCh <- ScanResult{
				CleanerName: clr.Name(),
				Category:    clr.Category(),
				Reclaimable: reclaimable,
				Err:         err,
			}
		}(c)
	}

	wg.Wait()
	close(resultsCh)

	var results []ScanResult
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}