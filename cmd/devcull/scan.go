package main

import (
	"fmt"
	"strings"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/engine"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [tool...]",
	Short: "Audit and estimate reclaimable disk space across installed tools",
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

		fmt.Println("🔍 Scanning developer caches...")
		fmt.Println()

		results := engine.Scan(activeCleaners)

		var totalReclaimable int64
		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("⚠️  %-12s error: %v\n", r.CleanerName, r.Err)
				continue
			}

			fmt.Printf("📦 %-12s %s\n", r.CleanerName, ui.FormatBytes(r.Reclaimable))
			totalReclaimable += r.Reclaimable
		}

		fmt.Printf("\n🎉 Estimated reclaimable space: %s\n", ui.FormatBytes(totalReclaimable))
		fmt.Println("Run 'devcull clean' to reclaim this space.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}