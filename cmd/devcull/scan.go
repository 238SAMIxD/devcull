package main

import (
	"fmt"
	"os"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [tool...]",
	Short: "Audit and estimate reclaimable disk space across installed tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		allCleaners := getAllCleaners()

		var activeCleaners []cleaner.Cleaner
		if len(args) > 0 {
			for _, c := range allCleaners {
				for _, arg := range args {
					if cleaner.MatchesArg(c, arg) {
						activeCleaners = append(activeCleaners, c)
						break
					}
				}
			}

			if len(activeCleaners) == 0 {
				fmt.Println("No matching tools or categories found for the provided arguments.")
				return nil
			}
		} else {
			activeCleaners = allCleaners
		}

		fmt.Println("🔍 Scanning developer caches...")

		results := engine.Scan(activeCleaners)

		grouped := make(map[cleaner.Category][]engine.ScanResult)
		var totalReclaimable int64
		var hasError bool

		for _, r := range results {
			if r.Skipped && len(args) == 0 {
				continue
			}
			grouped[r.Category] = append(grouped[r.Category], r)
			if r.Err == nil && !r.Skipped {
				totalReclaimable += r.Reclaimable
			}
			if r.Err != nil {
				hasError = true
			}
		}

		var displayCategories []cleaner.Category
		knownMap := make(map[cleaner.Category]bool)

		for _, cat := range cleaner.AllCategories() {
			if _, exists := grouped[cat]; exists {
				displayCategories = append(displayCategories, cat)
			}
			knownMap[cat] = true
		}

		var customCats []cleaner.Category
		for cat := range grouped {
			if !knownMap[cat] {
				customCats = append(customCats, cat)
			}
		}

		for i := 0; i < len(customCats)-1; i++ {
			for j := i + 1; j < len(customCats); j++ {
				if customCats[i] > customCats[j] {
					customCats[i], customCats[j] = customCats[j], customCats[i]
				}
			}
		}
		displayCategories = append(displayCategories, customCats...)

		for _, cat := range displayCategories {
			catResults, exists := grouped[cat]
			if !exists || len(catResults) == 0 {
				continue
			}

			var catTotal int64
			for _, r := range catResults {
				if r.Err == nil && !r.Skipped {
					catTotal += r.Reclaimable
				}
			}

			fmt.Printf("\n=== %s [%s] ===\n", cat, ui.FormatBytes(catTotal))

			for _, r := range catResults {
				if r.Skipped {
					fmt.Printf("⏭️  %-12s not installed\n", r.CleanerName)
					continue
				}
				if r.Err != nil {
					fmt.Printf("⚠️  %-12s error: %v\n", r.CleanerName, r.Err)
					continue
				}
				fmt.Printf("📦 %-12s %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimable))
			}
		}

		fmt.Printf("\n🎉 Total estimated reclaimable space: %s\n", ui.FormatBytes(totalReclaimable))
		fmt.Println("Run 'devcull clean' to reclaim this space.")

		if hasError {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
