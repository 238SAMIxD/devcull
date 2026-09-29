package main

import (
	"errors"
	"testing"

	"github.com/238SAMIxD/devcull/internal/engine"
)

func TestPrintCleanResults(t *testing.T) {
	results := []engine.Result{
		{CleanerName: "Node", Reclaimed: 1024, Skipped: false, Err: nil},
		{CleanerName: "FailTool", Reclaimed: 0, Skipped: false, Err: errors.New("some error")},
		{CleanerName: "SkippedTool", Reclaimed: 0, Skipped: true, Err: nil},
		{CleanerName: "PartialFail", Reclaimed: 512, Skipped: false, Err: errors.New("partial")},
	}

	total, hasErr, successfulRuns := printCleanResults(results, []string{}, false)

	if total != 1536 {
		t.Errorf("Expected total 1536, got %d", total)
	}
	if !hasErr {
		t.Error("Expected hasErr to be true")
	}
	if len(successfulRuns) != 2 {
		t.Errorf("Expected 2 successful runs, got %d", len(successfulRuns))
	}

	totalDry, hasErrDry, successfulRunsDry := printCleanResults(results, []string{}, true)
	if totalDry != 1536 {
		t.Errorf("Expected dry total 1536, got %d", totalDry)
	}
	if !hasErrDry {
		t.Error("Expected dry hasErr to be true")
	}
	if len(successfulRunsDry) != 0 {
		t.Errorf("Expected 0 successful runs in dry run, got %d", len(successfulRunsDry))
	}
}
