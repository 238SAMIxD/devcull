package main

import (
	"fmt"
	"os"

	"github.com/238SAMIxD/devcull/internal/cleaner"
	"github.com/238SAMIxD/devcull/internal/plugin"
)

func main() {
	nativeCleaners := cleaner.Native()
	pluginCleaners := plugin.LoadPlugins()
	allCleaners := append(nativeCleaners, pluginCleaners...)

	fmt.Printf("Successfully loaded %d cleaners!\n", len(allCleaners))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}