package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/serve"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm"},
	Short:   "Remove a project from registry",
	Long:    "Unregister a project from justvibin management. This removes the project from the registry and deletes the .justvibin marker file, but does not delete project files unless --files is specified.",
	Example: `justvibin remove myapp        # Remove specific project
justvibin rm myapp            # Alias for remove
justvibin remove myapp --files  # Also delete project files`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRemoveCmd,
}

func init() {
	removeCmd.Flags().Bool("files", false, "Also delete project files (dangerous!)")
	removeCmd.Flags().BoolP("yes", "y", false, "Skip confirmation for --files")
	rootCmd.AddCommand(removeCmd)
}

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	deleteFiles, _ := cmd.Flags().GetBool("files")
	skipConfirm, _ := cmd.Flags().GetBool("yes")

	var projectName string

	if len(args) > 0 {
		projectName = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("remove command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a justvibin project directory")
			logger.Info("Specify a project name: justvibin remove <name>")
			return errors.New("remove command failed")
		}
		marker, err := registry.ReadMarker(cwd)
		if err != nil {
			logger.Error("Failed to read project marker")
			return errors.New("remove command failed")
		}
		projectName = marker.Name
	}

	projectsPath, err := config.ProjectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return errors.New("remove command failed")
	}

	project, ok, err := registry.Get(projectsPath, projectName)
	if err != nil {
		logger.Error("Failed to load project registry")
		return errors.New("remove command failed")
	}
	if !ok {
		logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
		return errors.New("remove command failed")
	}

	projectPath := project.Path

	if deleteFiles && projectPath != "" {
		if !skipConfirm {
			fmt.Fprintf(cmd.OutOrStdout(), "DELETE all files in %s? [y/N] ", projectPath)
			reader := bufio.NewReader(cmd.InOrStdin())
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				logger.Info("Cancelled")
				return nil
			}
		}
	}

	if projectPath != "" {
		pidFile := filepath.Join(projectPath, serve.DefaultPIDFile)
		if pidData, err := os.ReadFile(pidFile); err == nil {
			if pid, err := strconv.Atoi(string(pidData)); err == nil {
				if process, err := os.FindProcess(pid); err == nil {
					_ = process.Signal(syscall.SIGTERM)
				}
			}
			_ = os.Remove(pidFile)
		}
	}

	removed, err := registry.Unregister(projectsPath, projectName)
	if err != nil {
		logger.Error("Failed to unregister project")
		return errors.New("remove command failed")
	}
	if !removed {
		logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
		return errors.New("remove command failed")
	}

	if deleteFiles && projectPath != "" {
		if err := os.RemoveAll(projectPath); err != nil {
			logger.Error(fmt.Sprintf("Failed to delete files: %v", err))
			return errors.New("remove command failed")
		}
		logger.Success(fmt.Sprintf("Deleted: %s", projectPath))
	} else if projectPath != "" {
		_ = os.Remove(filepath.Join(projectPath, ".justvibin"))
		logger.Info(fmt.Sprintf("Files remain at: %s", projectPath))
	}

	_ = console
	logger.Success(fmt.Sprintf("Removed: %s", projectName))
	return nil
}
