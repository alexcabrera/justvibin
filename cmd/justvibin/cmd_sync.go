package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexcabrera/justvibin/internal/config"
	execx "github.com/alexcabrera/justvibin/internal/exec"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/proxy"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [path]",
	Short: "Rebuild project registry by scanning for .justvibin files",
	Long:  "Scan directories for .justvibin marker files and rebuild the project registry. Useful when projects have been moved or the registry is corrupted. Use --clean to remove stale entries instead of scanning.",
	Example: `justvibin sync              # Scan home directory
justvibin sync ~/Code       # Scan specific path
justvibin sync --clean      # Remove entries for missing directories`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSyncCmd,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().Bool("clean", false, "Remove entries for missing directories instead of scanning")
}

type syncCommand struct {
	runner        execx.Runner
	projectsFile  func() (string, error)
	caddyfilePath func() (string, error)
	register      func(path, name string, port int, projectPath, template string) (registry.Project, error)
	generateCaddy func(context.Context, execx.Runner, string, string) error
	reloadProxy   func(context.Context, execx.Runner, string) error
	loadProjects  func(string) (map[string]registry.Project, error)
	unregister    func(string, string) (bool, error)
	saveProjects  func(string, map[string]registry.Project) error
}

var syncCommandFactory = defaultSyncCommand

func defaultSyncCommand() syncCommand {
	return syncCommand{
		runner:        execx.NewSystemRunner(),
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		register:      registry.Register,
		generateCaddy: proxy.GenerateCaddyfile,
		reloadProxy:   proxy.ReloadProxy,
		loadProjects:  registry.Load,
		unregister:    registry.Unregister,
		saveProjects:  registry.Save,
	}
}

func runSyncCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	cleanMode, _ := cmd.Flags().GetBool("clean")

	impl := syncCommandFactory()
	code := impl.run(context.Background(), args, console, logger, cleanMode)
	if code != 0 {
		return errors.New("sync command failed")
	}
	return nil
}

func (c syncCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, cleanMode bool) int {
	_ = console

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return 1
	}

	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}

	if cleanMode {
		return c.runClean(ctx, projectsPath, caddyfilePath, logger)
	}

	scanPath := os.Getenv("HOME")
	if len(args) > 0 {
		scanPath = args[0]
	}

	logger.Info(fmt.Sprintf("Scanning for .justvibin files in %s...", scanPath))

	if err := c.saveProjects(projectsPath, map[string]registry.Project{}); err != nil {
		logger.Error("Failed to clear registry")
		return 1
	}

	count := 0
	skipDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		".config":      true,
		"Library":      true,
		".Trash":       true,
	}

	err = filepath.WalkDir(scanPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != ".justvibin" {
			return nil
		}
		if count >= 100 {
			return filepath.SkipAll
		}

		projectDir := filepath.Dir(path)
		marker, err := registry.ReadMarker(projectDir)
		if err != nil {
			return nil
		}
		if marker.Name == "" || marker.Port == 0 {
			return nil
		}

		if c.register != nil {
			if _, err := c.register(projectsPath, marker.Name, marker.Port, projectDir, marker.Template); err != nil {
				return nil
			}
		}
		logger.Success(fmt.Sprintf("Found: %s (%s)", marker.Name, projectDir))
		count++
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		logger.Warn(fmt.Sprintf("Scan error: %v", err))
	}

	if c.generateCaddy != nil {
		_ = c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath)
	}
	if c.reloadProxy != nil {
		_ = c.reloadProxy(ctx, c.runner, caddyfilePath)
	}

	logger.Success(fmt.Sprintf("Synced %d project(s)", count))
	return 0
}

func (c syncCommand) runClean(ctx context.Context, projectsPath, caddyfilePath string, logger *logging.Logger) int {
	logger.Info("Cleaning stale registry entries...")

	projects, err := c.loadProjects(projectsPath)
	if err != nil {
		logger.Error("Failed to load projects")
		return 1
	}

	removed := 0
	for name, project := range projects {
		if _, err := os.Stat(project.Path); os.IsNotExist(err) {
			if c.unregister != nil {
				if _, err := c.unregister(projectsPath, name); err == nil {
					logger.Warn(fmt.Sprintf("Removing stale: %s (%s)", name, project.Path))
					removed++
				}
			}
			continue
		}
		if !registry.MarkerExists(project.Path) {
			if c.unregister != nil {
				if _, err := c.unregister(projectsPath, name); err == nil {
					logger.Warn(fmt.Sprintf("Removing stale: %s (%s)", name, project.Path))
					removed++
				}
			}
		}
	}

	if c.generateCaddy != nil {
		_ = c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath)
	}
	if c.reloadProxy != nil {
		_ = c.reloadProxy(ctx, c.runner, caddyfilePath)
	}

	logger.Success(fmt.Sprintf("Removed %d stale entries", removed))
	return 0
}
