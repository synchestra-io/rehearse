package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/rehearse/pkg/testscenario"
)

// ValidateCommand returns the "validate" cobra command.
func ValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate test scenarios and acceptance criteria",
		Long: `Validate the structural correctness of .test.md and .ac.md files
without executing any code. Checks include:

  - Scenario structure (title, description, step names, code annotations)
  - AC file structure (title, status, required sections, slug-filename match)
  - Cross-references (AC refs in scenarios resolve to actual files)
  - AC index sync (_acs/README.md matches files on disk)`,
		Args: cobra.MaximumNArgs(1),
		RunE: runValidate,
	}
	cmd.Flags().String("spec-root", "", "spec root directory for cross-reference resolution")
	cmd.Flags().Int("fail-fast", 0, "stop after N errors (default 1 when flag used without value)")
	cmd.Flags().Lookup("fail-fast").NoOptDefVal = "1"
	return cmd
}

func runValidate(cmd *cobra.Command, args []string) error {
	specRoot, _ := cmd.Flags().GetString("spec-root")
	if specRoot == "" {
		specRoot = "spec"
	}
	maxErrors, _ := cmd.Flags().GetInt("fail-fast")

	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return &CommandError{
			Err:  fmt.Errorf("path not found: %s", target),
			Code: 3,
		}
	}

	result, err := testscenario.ValidateAll(target, specRoot, maxErrors)
	if err != nil {
		return &CommandError{
			Err:  err,
			Code: 3,
		}
	}

	_, _ = fmt.Fprint(cmd.OutOrStdout(), testscenario.FormatValidationResult(result))

	if result.HasErrors() {
		return &CommandError{
			Err:  fmt.Errorf("validation failed with %d errors", len(result.Errors)),
			Code: 1,
		}
	}
	return nil
}
