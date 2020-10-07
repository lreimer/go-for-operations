package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var list *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ls",
	Short: "List directory contents",
	Long:  `List files and directories contained in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var lsCmd *exec.Cmd
		if *list {
			lsCmd = exec.Command("ls", "-l")
		} else {
			lsCmd = exec.Command("ls")
		}

		lsOut, err := lsCmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", string(lsOut))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	list = rootCmd.Flags().BoolP("long", "l", false, "List in long format.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ls" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ls")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
