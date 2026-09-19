package main

import (
	"fmt"
	"os"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
	"github.com/238SAMIxD/devcull/internal/stats"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var dryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean [tool...]",
	Short: "Run the cleaners (optionally specify which tools to clean)",
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

		state, err := stats.Load()
		if err != nil {
			return fmt.Errorf("failed to load stats: %w", err)
		}

		results := engine.Run(activeCleaners, dryRun)

		var sessionTotal int64
		hasArgs := len(args) > 0
		var hasError bool
		for _, r := range results {
			if r.Skipped {
				if hasArgs {
					fmt.Printf("⏭️  %s skipped (not installed)\n", r.CleanerName)
				}
				continue
			}
			if r.Err != nil {
				fmt.Printf("❌ %s failed: %v\n", r.CleanerName, r.Err)
				hasError = true
				continue
			}

			if dryRun {
				fmt.Printf("[DRY RUN] %s would reclaim %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
			} else {
				fmt.Printf("✅ %s reclaimed %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
				state.AddRun(r.CleanerName, r.Reclaimed)
			}
			sessionTotal += r.Reclaimed
		}

		if dryRun {
			fmt.Printf("\n🎉 Total space you COULD reclaim: %s\n", ui.FormatBytes(sessionTotal))
		} else {
			fmt.Printf("\n🎉 Total space reclaimed this session: %s\n", ui.FormatBytes(sessionTotal))
			if sessionTotal > 0 {
				if err := state.Save(); err != nil {
					fmt.Printf("⚠️ Failed to save stats: %v\n", err)
					hasError = true
				}
			}
		}

		if hasError {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate cleanup without deleting files")
	rootCmd.AddCommand(cleanCmd)
}
