package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/proxy"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Show or update project port",
	Long:  "Display or update the port for the current justvibin project. When updating, the registry and Caddy configuration are updated automatically.",
	Example: `justvibin port              # Show current port
justvibin port --set 8001   # Update port to 8001`,
	Args: cobra.NoArgs,
	RunE: runPortCmd,
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().Int("set", 0, "Set the port to use")
}

type portCommand struct {
	projectsFile     func() (string, error)
	caddyfilePath    func() (string, error)
	readMarker       func(string) (registry.Marker, error)
	updateMarkerPort func(string, int) (registry.Marker, error)
	updatePort       func(path, name string, port int) (registry.Project, error)
	generateCaddy    func(ctx context.Context, projectsPath, caddyfilePath string) error
	reloadProxy      func(ctx context.Context, caddyfilePath string) error
}

var portCommandFactory = defaultPortCommand

func defaultPortCommand() portCommand {
	return portCommand{
		projectsFile:     config.ProjectsFile,
		caddyfilePath:    config.CaddyfilePath,
		readMarker:       registry.ReadMarker,
		updateMarkerPort: registry.UpdateMarkerPort,
		updatePort:       registry.UpdatePort,
		generateCaddy: func(ctx context.Context, projectsPath, caddyfilePath string) error {
			return proxy.GenerateCaddyfile(ctx, nil, projectsPath, caddyfilePath)
		},
		reloadProxy: func(ctx context.Context, caddyfilePath string) error {
			return proxy.ReloadProxy(ctx, nil, caddyfilePath)
		},
	}
}

func runPortCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	setPort, _ := cmd.Flags().GetInt("set")

	impl := portCommandFactory()
	code := impl.run(context.Background(), console, logger, setPort)
	if code != 0 {
		return errors.New("port command failed")
	}
	return nil
}

func (c portCommand) run(ctx context.Context, console *ui.UI, logger *logging.Logger, setPort int) int {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Failed to get current directory")
		return 1
	}

	if !registry.MarkerExists(cwd) {
		logger.Error("Not a justvibin project directory")
		logger.Info("Run 'justvibin new' or 'justvibin register' first")
		return 1
	}

	marker, err := c.readMarker(cwd)
	if err != nil {
		logger.Error("Failed to read project marker")
		return 1
	}

	// If no --set flag, just show current port
	if setPort == 0 {
		fmt.Fprintf(os.Stdout, "%d\n", marker.Port)
		return 0
	}

	// Validate port
	if setPort < 1 || setPort > 65535 {
		logger.Error("Invalid port number (must be 1-65535)")
		return 1
	}

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return 1
	}

	// Update marker file
	if _, err := c.updateMarkerPort(cwd, setPort); err != nil {
		logger.Error(fmt.Sprintf("Failed to update marker: %v", err))
		return 1
	}

	// Update registry
	if _, err := c.updatePort(projectsPath, marker.Name, setPort); err != nil {
		logger.Error(fmt.Sprintf("Failed to update registry: %v", err))
		return 1
	}

	// Regenerate Caddyfile
	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Warn("Could not determine Caddyfile path")
	} else {
		if err := c.generateCaddy(ctx, projectsPath, caddyfilePath); err != nil {
			logger.Warn(fmt.Sprintf("Failed to regenerate Caddyfile: %v", err))
		} else {
			// Reload Caddy
			if err := c.reloadProxy(ctx, caddyfilePath); err != nil {
				logger.Warn(fmt.Sprintf("Failed to reload Caddy: %v", err))
			}
		}
	}

	logger.Success(fmt.Sprintf("Port updated to %d", setPort))
	return 0
}
