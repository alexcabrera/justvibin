package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var tunnelCmd = &cobra.Command{
	Use:   "tunnel [name]",
	Short: "Expose project via Cloudflare tunnel",
	Long:  "Create a quick public tunnel to expose your local project using cloudflared. Without arguments, tunnels the project in the current directory. Requires cloudflared to be installed.",
	Example: `justvibin tunnel              # Tunnel current project
justvibin tunnel myapp        # Tunnel specific project`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTunnelCmd,
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
}

func runTunnelCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	if _, err := exec.LookPath("cloudflared"); err != nil {
		logger.Error("cloudflared not installed")
		logger.Info("Install: brew install cloudflared")
		return errors.New("tunnel command failed")
	}

	var projectName string
	var port int

	if len(args) > 0 {
		projectName = args[0]
		projectsPath, err := config.ProjectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return errors.New("tunnel command failed")
		}
		project, ok, err := registry.Get(projectsPath, projectName)
		if err != nil {
			logger.Error("Failed to load project registry")
			return errors.New("tunnel command failed")
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
			return errors.New("tunnel command failed")
		}
		port = project.Port
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("tunnel command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a justvibin project directory")
			logger.Info("Run 'justvibin new' or 'justvibin register' first")
			return errors.New("tunnel command failed")
		}
		marker, err := registry.ReadMarker(cwd)
		if err != nil {
			logger.Error("Failed to read project marker")
			return errors.New("tunnel command failed")
		}
		projectName = marker.Name
		port = marker.Port
	}

	if !isPortInUse(port) {
		logger.Error("Project not running")
		logger.Info("Start first: justvibin start")
		return errors.New("tunnel command failed")
	}

	_ = console
	logger.Info(fmt.Sprintf("Starting quick tunnel for %s...", projectName))
	logger.Info("Press Ctrl+C to stop")

	tunnelCmd := exec.Command("cloudflared", "tunnel", "--url", fmt.Sprintf("https://localhost:%d", port))
	tunnelCmd.Stdout = os.Stdout
	tunnelCmd.Stderr = os.Stderr
	tunnelCmd.Stdin = os.Stdin
	tunnelCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: false}

	if err := tunnelCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return nil
		}
		logger.Error(fmt.Sprintf("Tunnel failed: %v", err))
		return errors.New("tunnel command failed")
	}
	return nil
}
