package main

import (
	"fmt"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/spf13/cobra"
)

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "devcull",
	Short: "A blazing-fast CLI to reclaim disk space from developer tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &cleaner.PnpmCleaner{}

		if !p.IsInstalled() {
			fmt.Printf("%s not installed, skipping.\n", p.Name())
			return nil
		}

		reclaimed, err := p.Clean(dryRun)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("[DRY RUN] %s would reclaim %d bytes\n", p.Name(), reclaimed)
		} else {
			fmt.Printf("%s reclaimed %d bytes\n", p.Name(), reclaimed)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate cleanup without deleting files")
}