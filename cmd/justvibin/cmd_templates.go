package main

import (
	"errors"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List installed templates",
	Long:  "List template plugins installed in ~/.config/justvibin/templates/. Shows template name, description, serve type, and source URL for each entry. Use --json for machine-readable output in scripts and CI.",
	Example: "justvibin templates\njustvibin templates --json\njustvibin --no-color templates\njustvibin --quiet templates",
	RunE:  runTemplatesCmd,
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}

func runTemplatesCmd(cmd *cobra.Command, _ []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	jsonOutput := output.JSON

	templatesDir, err := config.TemplatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return errors.New("templates command failed")
	}
	installed, err := loadInstalledTemplates(templatesDir)
	if err != nil {
		logger.Error("Failed to load installed templates")
		return errors.New("templates command failed")
	}

	if len(installed) == 0 {
		logger.Info("No templates installed.")
		logger.Info("Install with: justvibin install <git-url>")
		logger.Info("Or see official: justvibin install --list-official")
		return nil
	}

	if jsonOutput {
		payload, err := templatesJSON(installed)
		if err != nil {
			logger.Error("Failed to render JSON")
			return errors.New("templates command failed")
		}
		console.PrintHelp(payload)
		return nil
	}

	console.PrintHelp(templatesText(installed, output.Styled))
	return nil
}
