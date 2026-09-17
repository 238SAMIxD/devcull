package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcull",
	Short: "A blazing-fast CLI to reclaim disk space from developer tools",
}

func Execute() error {
	return rootCmd.Execute()
}