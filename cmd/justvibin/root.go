package main

import (
	"context"
	"os"

	"github.com/alexcabrera/justvibin/internal/version"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	flagQuiet   bool
	flagVerbose bool
	flagNoColor bool
	flagJSON    bool
)

var rootCmd = &cobra.Command{
	Use:     "justvibin",
	Short:   "CLI for scaffolding web application projects",
	Long:    "justvibin scaffolds web projects from curated templates and manages a local HTTPS proxy. Use it to create new projects, install template plugins, and run setup tasks. Global flags let you control output formatting for scripting and CI workflows.",
	Example: "justvibin new myapp\njustvibin templates --json\njustvibin setup --check\njustvibin --help",
	Version: version.Version,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress non-essential output. Default: false")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output. Default: false")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output. Default: false")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format where supported. Default: false")
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose")
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
