package cli

import (
	"os"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

func Execute() {
	//nolint: exhaustruct
	root := &cobra.Command{
		Use:   "monitor",
		Short: "project of software engineering course",
	}

	root.AddCommand(serve.New())

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
