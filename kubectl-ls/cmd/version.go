package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

// SetVersion set the application version for consumption in the output of the command.
func SetVersion(v string) {
	version = v
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  "Display version number and commit hash of application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kubectl-ls %v\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
