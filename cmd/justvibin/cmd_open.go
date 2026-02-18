package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [name]",
	Short: "Open project in browser",
	Long:  "Open the project URL in your default browser. Without arguments, opens the project in the current directory.",
	Example: `justvibin open              # Open current project
justvibin open myapp        # Open specific project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runOpenCmd,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpenCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	var projectName string

	if len(args) > 0 {
		projectName = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("open command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a justvibin project directory")
			logger.Info("Run 'justvibin new' or 'justvibin register' first")
			return errors.New("open command failed")
		}
		marker, err := registry.ReadMarker(cwd)
		if err != nil {
			logger.Error("Failed to read project marker")
			return errors.New("open command failed")
		}
		projectName = marker.Name
	}

	projectsPath, err := config.ProjectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return errors.New("open command failed")
	}

	_, ok, err := registry.Get(projectsPath, projectName)
	if err != nil {
		logger.Error("Failed to load project registry")
		return errors.New("open command failed")
	}
	if !ok {
		logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
		return errors.New("open command failed")
	}

	url := fmt.Sprintf("https://%s.localhost", projectName)
	if err := openBrowser(url); err != nil {
		logger.Error(fmt.Sprintf("Failed to open browser: %v", err))
		return errors.New("open command failed")
	}

	_ = console
	logger.Success(fmt.Sprintf("Opened: %s", url))
	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
