package main

import (
	"fmt"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "devcull",
	Short: "A blazing-fast CLI to reclaim disk space from developer tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		cleaners := []cleaner.Cleaner{
			&cleaner.PnpmCleaner{},
			// Future tools go here:
			// &cleaner.BrewCleaner{},
			// &cleaner.UvCleaner{},
		}

		results := engine.Run(cleaners, dryRun)

		var totalReclaimed int64
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("❌ %s failed: %v\n", r.CleanerName, r.Err)
				continue
			}
			
			if dryRun {
				fmt.Printf("[DRY RUN] %s would reclaim %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
			} else {
				fmt.Printf("✅ %s reclaimed %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
			}
			totalReclaimed += r.Reclaimed
		}

		fmt.Printf("\n🎉 Total space reclaimed: %s\n", ui.FormatBytes(totalReclaimed))
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate cleanup without deleting files")
}