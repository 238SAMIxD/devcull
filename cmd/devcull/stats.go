package main

import (
	"fmt"

	"github.com/238SAMIxD/devcull/internal/stats"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your all-time disk space reclamation stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := stats.Load()
		if err != nil {
			return fmt.Errorf("failed to load stats: %w", err)
		}

		if state.AllTimeTotal == 0 {
			fmt.Println("No stats yet! Run 'devcull clean' to start reclaiming space.")
			return nil
		}

		fmt.Printf("🏆 All-Time Total Reclaimed: %s\n\n", ui.FormatBytes(state.AllTimeTotal))
		
		for tool, ts := range state.Tools {
			fmt.Printf("📦 %s\n", tool)
			fmt.Printf("   Total: %s\n", ui.FormatBytes(ts.TotalReclaimed))
			fmt.Printf("   Recent Runs:\n")
			for i, run := range ts.RecentRuns {
				fmt.Printf("     %d. %s\n", i+1, ui.FormatBytes(run))
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}