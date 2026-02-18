package main

import (
	"errors"

	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <name>",
	Short: "Remove an installed template plugin",
	Long:  "Remove an installed template plugin by name. Use --all to remove all installed templates with confirmation. If projects are using the template, the command will warn and ask before proceeding.",
	Example: "justvibin uninstall hypertext\njustvibin uninstall --all\njustvibin uninstall hypertext --quiet\njustvibin --no-color uninstall hypertext",
	Args:   cobra.MaximumNArgs(1),
	RunE:   runUninstallCmd,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().Bool("all", false, "Remove all installed templates with confirmation. Default: false")
}

func runUninstallCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	cmdImpl := uninstallCommandFactory()
	cmdImpl.confirm = newUninstallConfirmer(cmd, logger)

	argsToRun := make([]string, 0, 2)
	if all {
		argsToRun = append(argsToRun, "--all")
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger)
	if code != 0 {
		return errors.New("uninstall command failed")
	}
	return nil
}

func newUninstallConfirmer(cmd *cobra.Command, logger *logging.Logger) func(string) (bool, error) {
	return func(question string) (bool, error) {
		if output := getOutputSettings(cmd); output.Quiet {
			return false, nil
		}
		confirmed := false
		prompt := huh.NewConfirm().Title(question).Value(&confirmed)
		if err := prompt.Run(); err != nil {
			logger.Error("Failed to read confirmation")
			return false, err
		}
		return confirmed, nil
	}
}
