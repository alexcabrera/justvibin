package main

import (
	"errors"

	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <template-name>",
	Short: "Update installed templates from their source repositories",
	Long:  "Re-fetch a template from its original git source URL. Use --all to update all templates at once. Requires git to be installed and the template to have a valid .source file.",
	Example: `justvibin update hypertext    # Update specific template
justvibin update --all        # Update all installed templates`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdateCmd,
}

func init() {
	updateCmd.Flags().Bool("all", false, "Update all installed templates")
	rootCmd.AddCommand(updateCmd)
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	updateAll, _ := cmd.Flags().GetBool("all")

	if !updateAll && len(args) == 0 {
		logger.Error("Missing template name")
		logger.Info("Usage: justvibin update <template-name>")
		logger.Info("Or: justvibin update --all")
		return errors.New("update command failed")
	}

	cmdImpl := updateCommandFactory()

	argsToRun := make([]string, 0, 2)
	if updateAll {
		argsToRun = append(argsToRun, "--all")
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger, output.Styled)
	if code != 0 {
		return errors.New("update command failed")
	}
	_ = console
	return nil
}


