package cli

import (
	"context"
	"errors"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Run executes the rehearse CLI with the given arguments.
func Run(
	args []string,
	fatal func(error),
) {
	rootCmd := &cobra.Command{
		Use:           "rehearse",
		Short:         "Markdown-native test framework",
		Long:          "Rehearse turns markdown specifications into executable test scenarios.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetErr(os.Stderr)

	rootCmd.AddCommand(
		RunCommand(),
		ListCommand(),
		versionCommand(),
	)

	rootCmd.SetArgs(args[1:])
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		fatal(err)
	}
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Printf("rehearse %s (%s) %s\n", version, commit, date)
		},
	}
}

// Command returns the "test" cobra command for embedding in a parent CLI (e.g., synchestra test).
// It wraps RunCommand and ListCommand under a "test" parent.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run and manage test scenarios",
	}
	cmd.AddCommand(
		RunCommand(),
		ListCommand(),
	)
	return cmd
}

// CommandError wraps an error with an exit code.
type CommandError struct {
	Err  error
	Code int
}

func (e *CommandError) Error() string { return e.Err.Error() }
func (e *CommandError) ExitCode() int { return e.Code }
func (e *CommandError) Unwrap() error { return e.Err }

// NewCommandError creates a CommandError with exit code 1.
func NewCommandError(err error) error {
	return &CommandError{Err: err, Code: 1}
}

// IsCommandError returns true if err is or wraps a CommandError.
func IsCommandError(err error) bool {
	var ce *CommandError
	return errors.As(err, &ce)
}
