package main

import (
	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/plugin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcull",
	Short: "A blazing-fast CLI to reclaim disk space from developer tools",
}

func Execute() error {
	return rootCmd.Execute()
}

func getAllCleaners() []cleaner.Cleaner {
	all := cleaner.Native()
	all = append(all, plugin.LoadPlugins()...)
	return all
}