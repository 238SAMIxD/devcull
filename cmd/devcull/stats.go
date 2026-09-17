package main

import (
	"fmt"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/stats"
	"github.com/238SAMIxD/devcull/internal/ui"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View historical cleanup statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := stats.Load()
		if err != nil {
			return fmt.Errorf("failed to load stats: %w", err)
		}

		if s.AllTimeTotal == 0 {
			fmt.Println("No cleanup history found. Run 'devcull clean' to get started!")
			return nil
		}

		fmt.Println("📊 Devcull Historical Stats")
		fmt.Println("===========================")
		fmt.Printf("Total Space Reclaimed: %s\n\n", ui.FormatBytes(s.AllTimeTotal))

		toolToCategory := make(map[string]cleaner.Category)
		for _, c := range cleaner.All() {
			toolToCategory[c.Name()] = c.Category()
		}

		groupedStats := make(map[cleaner.Category]map[string]int64)
		categoryTotals := make(map[cleaner.Category]int64)

		for toolName, toolStat := range s.Tools {
			if toolStat == nil || toolStat.TotalReclaimed == 0 {
				continue
			}

			cat, exists := toolToCategory[toolName]
			if !exists {
				cat = cleaner.CategorySystem
			}

			if groupedStats[cat] == nil {
				groupedStats[cat] = make(map[string]int64)
			}
			
			groupedStats[cat][toolName] = toolStat.TotalReclaimed
			categoryTotals[cat] += toolStat.TotalReclaimed
		}

		fmt.Println("🏆 All-Time Leaderboard")

		for _, cat := range cleaner.AllCategories() {
			toolsInCat, exists := groupedStats[cat]
			if !exists || len(toolsInCat) == 0 {
				continue
			}

			fmt.Printf("\n=== %s [%s] ===\n", cat, ui.FormatBytes(categoryTotals[cat]))
			for toolName, reclaimed := range toolsInCat {
				fmt.Printf("  %-12s %s\n", toolName, ui.FormatBytes(reclaimed))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}