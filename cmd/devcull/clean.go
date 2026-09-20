package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
	"github.com/238SAMIxD/devcull/internal/stats"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var dryRun bool
var yesRun bool

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

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if !dryRun && !yesRun {
			fmt.Println("Scanning...")
			scanResults := engine.Scan(ctx, activeCleaners)
			var totalReclaimable int64
			var scanErr error
			for _, r := range scanResults {
				if r.Err != nil {
					scanErr = r.Err
					break
				}
				if !r.Skipped {
					totalReclaimable += r.Reclaimable
				}
			}

			if scanErr != nil {
				fmt.Printf("❌ Preflight scan failed: %v\n", scanErr)
				os.Exit(1)
			}

			if totalReclaimable == 0 {
				fmt.Println("Nothing to clean.")
				return nil
			}

			fmt.Printf("Are you sure you want to reclaim %s? [y/N] ", ui.FormatBytes(totalReclaimable))
			reader := bufio.NewReader(os.Stdin)
			resp, _ := reader.ReadString('\n')
			resp = strings.ToLower(strings.TrimSpace(resp))
			if resp != "y" && resp != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		results := engine.Run(ctx, activeCleaners, dryRun)

		var sessionTotal int64
		hasArgs := len(args) > 0
		var hasError bool

		type runStat struct {
			name      string
			reclaimed int64
		}
		var successfulRuns []runStat

		for _, r := range results {
			if r.Skipped {
				if hasArgs {
					fmt.Printf("⏭️  %s skipped (not installed)\n", r.CleanerName)
				}
				continue
			}
			if r.Reclaimed > 0 {
				sessionTotal += r.Reclaimed
				if !dryRun && r.Err != nil {
					successfulRuns = append(successfulRuns, runStat{name: r.CleanerName, reclaimed: r.Reclaimed})
				}
			}

			if r.Err != nil {
				fmt.Printf("❌ %s failed: %v\n", r.CleanerName, r.Err)
				if r.Reclaimed > 0 {
					if dryRun {
						fmt.Printf("   (Partially would reclaim %s)\n", ui.FormatBytes(r.Reclaimed))
					} else {
						fmt.Printf("   (Partially reclaimed %s)\n", ui.FormatBytes(r.Reclaimed))
					}
				}
				hasError = true
				continue
			}

			if r.Reclaimed > 0 {
				if dryRun {
					fmt.Printf("[DRY RUN] %s would reclaim %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
				} else {
					fmt.Printf("✅ %s reclaimed %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimed))
					successfulRuns = append(successfulRuns, runStat{name: r.CleanerName, reclaimed: r.Reclaimed})
				}
			}
		}

		if dryRun {
			fmt.Printf("\n🎉 Total space you COULD reclaim: %s\n", ui.FormatBytes(sessionTotal))
		} else {
			fmt.Printf("\n🎉 Total space reclaimed this session: %s\n", ui.FormatBytes(sessionTotal))
			if sessionTotal > 0 {
				state, err := stats.LoadAndLock()
				if err != nil {
					fmt.Printf("⚠️ Failed to load stats for saving: %v\n", err)
					hasError = true
				} else {
					defer state.Unlock()
					for _, run := range successfulRuns {
						state.AddRun(run.name, run.reclaimed)
					}
					if err := state.Save(); err != nil {
						fmt.Printf("⚠️ Failed to save stats: %v\n", err)
						hasError = true
					}
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
	cleanCmd.Flags().BoolVarP(&yesRun, "yes", "y", false, "Skip confirmation prompt")
	rootCmd.AddCommand(cleanCmd)
}
