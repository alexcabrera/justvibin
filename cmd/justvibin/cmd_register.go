package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alexcabrera/justvibin/internal/config"
	execx "github.com/alexcabrera/justvibin/internal/exec"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/proxy"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register [name]",
	Short: "Register existing directory with justvibin",
	Long:  "Register an existing project directory with justvibin for serving and management. Automatically detects project type, assigns a port, creates .justvibin marker, and updates proxy configuration.",
	Example: `justvibin register              # Register current dir with auto-detected name
justvibin register myapp        # Register with specific name
justvibin register -t static    # Specify template type`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRegisterCmd,
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("template", "t", "", "Specify template type (static, django, node)")
}

var registerNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

type registerCommand struct {
	runner        execx.Runner
	projectsFile  func() (string, error)
	caddyfilePath func() (string, error)
	nextPort      func(string) (int, error)
	register      func(path, name string, port int, projectPath, template string) (registry.Project, error)
	writeMarker   func(projectDir, name, template string, port int) (registry.Marker, error)
	generateCaddy func(context.Context, execx.Runner, string, string) error
	reloadProxy   func(context.Context, execx.Runner, string) error
	markerExists  func(string) bool
}

var registerCommandFactory = defaultRegisterCommand

func defaultRegisterCommand() registerCommand {
	return registerCommand{
		runner:        execx.NewSystemRunner(),
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		nextPort:      registry.NextPort,
		register:      registry.Register,
		writeMarker:   registry.WriteMarker,
		generateCaddy: proxy.GenerateCaddyfile,
		reloadProxy:   proxy.ReloadProxy,
		markerExists:  registry.MarkerExists,
	}
}

func runRegisterCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	templateName, _ := cmd.Flags().GetString("template")

	impl := registerCommandFactory()
	code := impl.run(context.Background(), args, console, logger, templateName)
	if code != 0 {
		return errors.New("register command failed")
	}
	return nil
}

func (c registerCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, templateName string) int {
	_ = console

	projectDir, err := os.Getwd()
	if err != nil {
		logger.Error("Failed to get current directory")
		return 1
	}

	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	} else {
		projectName = filepath.Base(projectDir)
	}

	if !registerNamePattern.MatchString(projectName) {
		logger.Error("Invalid name: use letters, numbers, hyphens, underscores. Must start with a letter.")
		return 1
	}

	if c.markerExists(projectDir) {
		marker, err := registry.ReadMarker(projectDir)
		if err == nil {
			logger.Warn(fmt.Sprintf("Directory already registered as '%s'", marker.Name))
			return 0
		}
	}

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return 1
	}

	existing, ok, _ := registry.Get(projectsPath, projectName)
	if ok && existing.Path != projectDir {
		logger.Error(fmt.Sprintf("Name '%s' already used by %s", projectName, existing.Path))
		return 1
	}

	if templateName == "" {
		templateName = detectProjectType(projectDir)
		if templateName == "" {
			templateName = "static"
		}
		logger.Info(fmt.Sprintf("Detected template type: %s", templateName))
	}

	port, err := c.nextPort(projectsPath)
	if err != nil {
		logger.Error("Failed to assign port")
		return 1
	}

	if c.writeMarker != nil {
		if _, err := c.writeMarker(projectDir, projectName, templateName, port); err != nil {
			logger.Error("Failed to write .justvibin marker")
			return 1
		}
	}

	if c.register != nil {
		if _, err := c.register(projectsPath, projectName, port, projectDir, templateName); err != nil {
			logger.Error("Failed to register project")
			return 1
		}
	}

	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}

	if c.generateCaddy != nil {
		if err := c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath); err != nil {
			logger.Error("Failed to generate Caddyfile")
			return 1
		}
	}

	if c.reloadProxy != nil {
		if err := c.reloadProxy(ctx, c.runner, caddyfilePath); err != nil {
			logger.Error("Failed to reload proxy")
			return 1
		}
	}

	logger.Success(fmt.Sprintf("Registered: %s", projectName))
	logger.Info(fmt.Sprintf("URL: https://%s.localhost", projectName))
	logger.Info("Start: justvibin start")
	return 0
}

func detectProjectType(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, "manage.py")); err == nil {
		return "django-hypermedia"
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		return "node"
	}
	if _, err := os.Stat(filepath.Join(dir, "index.html")); err == nil {
		return "static"
	}
	return ""
}
