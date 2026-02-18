package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alexcabrera/justvibin/internal/config"
	"github.com/alexcabrera/justvibin/internal/logging"
	"github.com/alexcabrera/justvibin/internal/manifest"
	"github.com/alexcabrera/justvibin/internal/registry"
	"github.com/alexcabrera/justvibin/internal/serve"
	"github.com/alexcabrera/justvibin/internal/ui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a project server",
	Long:  "Start a development or production server for a justvibin project. Without arguments, starts the project in the current directory. Supports static file serving and command-based servers defined in the template manifest.",
	Example: `justvibin start              # Start current project
justvibin start myapp        # Start specific project
justvibin start --prod       # Start in production mode`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Bool("prod", false, "Run in production mode")
}

type startCommand struct {
	templatesDir  func() (string, error)
	projectsFile  func() (string, error)
	readMarker    func(string) (registry.Marker, error)
	readFile      func(string) ([]byte, error)
	startStatic   func(ctx context.Context, runner serve.CommandRunner, port int, root string) (int, error)
	startCommand  func(ctx context.Context, dir string, cmd string, port int, portEnv string) (int, error)
	isPortInUse   func(int) bool
}

var startCommandFactory = defaultStartCommand

func defaultStartCommand() startCommand {
	return startCommand{
		templatesDir:  config.TemplatesDir,
		projectsFile:  config.ProjectsFile,
		readMarker:    registry.ReadMarker,
		readFile:      os.ReadFile,
		startStatic:   serve.StartStaticServer,
		startCommand:  startCommandServer,
		isPortInUse:   isPortInUse,
	}
}

func runStartCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	prodMode, _ := cmd.Flags().GetBool("prod")

	impl := startCommandFactory()
	code := impl.run(context.Background(), args, console, logger, prodMode)
	if code != 0 {
		return errors.New("start command failed")
	}
	return nil
}

func (c startCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, prodMode bool) int {
	var projectDir string
	var projectName string

	if len(args) > 0 {
		projectName = args[0]
		projectsPath, err := c.projectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return 1
		}
		project, ok, err := registry.Get(projectsPath, projectName)
		if err != nil {
			logger.Error("Failed to load project registry")
			return 1
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
			return 1
		}
		projectDir = project.Path
	} else {
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
		projectDir = cwd
	}

	marker, err := c.readMarker(projectDir)
	if err != nil {
		logger.Error("Failed to read project marker")
		return 1
	}
	projectName = marker.Name
	port := marker.Port
	templateName := marker.Template

	if c.isPortInUse(port) {
		logger.Warn(fmt.Sprintf("Project already running on port %d", port))
		logger.Info(fmt.Sprintf("URL: https://%s.localhost", projectName))
		return 0
	}

	templatesDir, err := c.templatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}

	templateDir := filepath.Join(templatesDir, templateName)
	manifestPath := filepath.Join(templateDir, "justvibin.toml")

	serveType := "static"
	var mf manifest.Manifest
	if data, err := c.readFile(manifestPath); err == nil {
		if parsed, err := manifest.Parse(data); err == nil {
			mf = parsed
			if parsed.Serve.Type != "" {
				serveType = parsed.Serve.Type
			}
		}
	}

	logger.Info(fmt.Sprintf("Starting %s on port %d...", projectName, port))

	switch serveType {
	case "static":
		staticRoot := projectDir
		if mf.Serve.Static.Root != "" {
			staticRoot = filepath.Join(projectDir, mf.Serve.Static.Root)
		}
		_, err := c.startStatic(ctx, nil, port, staticRoot)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start static server: %v", err))
			return 1
		}
	case "command":
		cmdStr := manifest.ServeCommand(mf, modeString(prodMode))
		if cmdStr == "" {
			logger.Error("No serve command defined in manifest")
			return 1
		}
		portEnv := mf.Serve.PortEnv
		if portEnv == "" {
			portEnv = "PORT"
		}
		_, err := c.startCommand(ctx, projectDir, cmdStr, port, portEnv)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
			return 1
		}
	default:
		logger.Error(fmt.Sprintf("Unknown serve type: %s", serveType))
		return 1
	}

	logger.Success(fmt.Sprintf("Started: https://%s.localhost", projectName))
	return 0
}

func modeString(prod bool) string {
	if prod {
		return "prod"
	}
	return "dev"
}

func startCommandServer(ctx context.Context, dir string, cmdStr string, port int, portEnv string) (int, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%d", portEnv, port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid := cmd.Process.Pid
	pidFile := filepath.Join(dir, serve.DefaultPIDFile)
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		_ = cmd.Process.Kill()
		return 0, err
	}

	return pid, nil
}
