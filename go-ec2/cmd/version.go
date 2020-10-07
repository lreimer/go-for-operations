package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string
var commit string

// SetVersion set the application version for consumption in the output of the command.
func SetVersion(v string) {
	version = v
}

// SetCommit set the application commit for consumption in the output of the command.
func SetCommit(c string) {
	commit = c
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  "Display information on version number and Git commit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-ec2 %v %v\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
