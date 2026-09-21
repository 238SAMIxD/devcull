package main

import (
	"errors"
	"testing"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
)

func TestPrintScanResults(t *testing.T) {
	results := []engine.ScanResult{
		{CleanerName: "Node", Category: cleaner.CategoryNode, Reclaimable: 1024, Skipped: false, Err: nil},
		{CleanerName: "Python", Category: cleaner.CategoryPython, Reclaimable: 0, Skipped: false, Err: errors.New("fail")},
		{CleanerName: "Java", Category: cleaner.CategoryJava, Reclaimable: 0, Skipped: true, Err: nil},
		{CleanerName: "Go", Category: cleaner.CategoryGo, Reclaimable: 2048, Skipped: false, Err: nil},
	}

	total, hasErr := printScanResults(results, []string{})

	if total != 3072 {
		t.Errorf("Expected total 3072, got %d", total)
	}
	if !hasErr {
		t.Error("Expected hasErr to be true")
	}
}
