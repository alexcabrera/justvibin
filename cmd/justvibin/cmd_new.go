package main

import (
	"context"
	"errors"
	"os"

	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new project from a template",
	Long:  "Create a new project directory from a curated template. Templates are cloned from git repositories and initialized with a fresh git repo. If you omit the name, you can pass it via --name or be prompted when running interactively. For headless usage, provide --name or a positional name along with any flags.",
	Example: "justvibin new myapp\njustvibin new --template hypertext myapp\njustvibin new --local ./templates/hypertext --name myapp\njustvibin --json templates",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runNewCmd,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("template", "t", "django-hypermedia", "Specify template name. Default: django-hypermedia")
	newCmd.Flags().String("local", "", "Use local template directory instead of cloning. Default: empty")
	newCmd.Flags().StringP("name", "n", "", "Project name (alternative to positional arg). Default: empty")
}

func runNewCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	if name == "" && len(args) > 0 {
		name = args[0]
	}

	templateName, err := cmd.Flags().GetString("template")
	if err != nil {
		return err
	}
	templateChanged := cmd.Flags().Changed("template")

	localPath, err := cmd.Flags().GetString("local")
	if err != nil {
		return err
	}

	newArgs := make([]string, 0, 6)
	if name != "" {
		newArgs = append(newArgs, name)
	}
	if localPath != "" {
		newArgs = append(newArgs, "--local", localPath)
	}
	if templateChanged {
		newArgs = append(newArgs, "--template", templateName)
	}

	cmdImpl := newCommandFactory()
	interactive := term.IsTerminal(int(os.Stdin.Fd()))
	code := cmdImpl.run(context.Background(), newArgs, console, logger, interactive)
	if code != 0 {
		return errors.New("new command failed")
	}
	return nil
}
