package main

import (
	"errors"

	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <git-url>",
	Short: "Install a template plugin from a git repository",
	Long:  "Clone a template repository and install it as a plugin. Templates must contain a justvibin.toml manifest file with template metadata. Use --list-official to browse curated templates without installing, or --name to override the installed template name.",
	Example: `justvibin install https://github.com/acme/my-template.git
justvibin install --name custom-name https://github.com/acme/my-template.git
justvibin install --list-official
justvibin --quiet install --list-official`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInstallCmd,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("name", "n", "", "Override template name (default: from manifest). Default: empty")
	installCmd.Flags().Bool("list-official", false, "List official templates instead of installing. Default: false")
}

func runInstallCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	listOfficial, err := cmd.Flags().GetBool("list-official")
	if err != nil {
		return err
	}

	if !listOfficial && len(args) == 0 {
		logger.Error("Missing template git URL")
		return errors.New("install command failed")
	}

	cmdImpl := installCommandFactory()

	argsToRun := make([]string, 0, 4)
	if listOfficial {
		argsToRun = append(argsToRun, "--list-official")
	}
	if name != "" {
		argsToRun = append(argsToRun, "--name", name)
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger, output.Styled)
	if code != 0 {
		return errors.New("install command failed")
	}
	return nil
}
