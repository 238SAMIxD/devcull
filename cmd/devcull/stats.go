package main

import (
	"fmt"

	"sort"

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
		for _, c := range getAllCleaners() {
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

		type statEntry struct {
			name      string
			reclaimed int64
		}

		var activeCategories []cleaner.Category

		for _, cat := range cleaner.AllCategories() {
			if _, exists := groupedStats[cat]; exists {
				activeCategories = append(activeCategories, cat)
			}
		}

		var extraCategories []cleaner.Category
		for cat := range groupedStats {
			isBuiltIn := false
			for _, builtIn := range cleaner.AllCategories() {
				if cat == builtIn {
					isBuiltIn = true
					break
				}
			}
			if !isBuiltIn {
				extraCategories = append(extraCategories, cat)
			}
		}

		sort.Slice(extraCategories, func(i, j int) bool {
			return extraCategories[i] < extraCategories[j]
		})

		activeCategories = append(activeCategories, extraCategories...)

		for _, cat := range activeCategories {
			toolsInCat, exists := groupedStats[cat]
			if !exists || len(toolsInCat) == 0 {
				continue
			}

			fmt.Printf("\n=== %s [%s] ===\n", cat, ui.FormatBytes(categoryTotals[cat]))

			var entries []statEntry
			for name, reclaimed := range toolsInCat {
				entries = append(entries, statEntry{name, reclaimed})
			}

			sort.Slice(entries, func(i, j int) bool {
				if entries[i].reclaimed == entries[j].reclaimed {
					return entries[i].name < entries[j].name
				}
				return entries[i].reclaimed > entries[j].reclaimed
			})

			for _, entry := range entries {
				fmt.Printf("  %-12s %s\n", entry.name, ui.FormatBytes(entry.reclaimed))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
