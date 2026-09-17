package main

import (
	"fmt"
	"strings"

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
		allCleaners := cleaner.All()
		
		var activeCleaners []cleaner.Cleaner
		if len(args) > 0 {
			requested := make(map[string]bool)
			for _, arg := range args {
				requested[cleaner.ResolveAlias(arg)] = true
			}

			for _, c := range allCleaners {
				if requested[strings.ToLower(c.Name())] {
					activeCleaners = append(activeCleaners, c)
				}
			}

			if len(activeCleaners) == 0 {
				fmt.Println("No matching tools found for the provided arguments.")
				return nil
			}
		} else {
			activeCleaners = allCleaners
		}

		results := engine.Run(activeCleaners, dryRun)

		state, err := stats.Load()
		if err != nil {
			return fmt.Errorf("failed to load stats: %w", err)
		}

		var sessionTotal int64
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("❌ %s failed: %v\n", r.CleanerName, r.Err)
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
				}
			}
		}

		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate cleanup without deleting files")
	rootCmd.AddCommand(cleanCmd)
}