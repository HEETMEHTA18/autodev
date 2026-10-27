package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print AutoDev version and build information",
		Long:  `Print the AutoDev CLI version along with build and runtime details. Use --json for machine-readable output.`,
		Example: `  autodev version
  autodev version --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := map[string]string{
				"version":    rootCmd.Version,
				"go_version": runtime.Version(),
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
			}
			if jsonOut {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}
			fmt.Printf("  autodev version %s\n", styleBold.Render(info["version"]))
			fmt.Printf("  %s %s %s/%s\n", styleDim.Render("build"), info["go_version"], info["os"], info["arch"])
			return nil
		},
	}
}
