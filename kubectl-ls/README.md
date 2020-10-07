# kubectl-ls Plugin using Cobra

This demo shows how to create a simple Cobra CLI app to display and output
a directory listing similar to `ls`. By following some naming conventions the
application can be use as a `kubectl` plugin. Example usage:
```bash
$ kubectl ls
$ kubectl ls -l
$ kubectl ls --help
$ kubectl ls help
$ kubectl ls version
```

## Initial application creation

The initial skaffolding of the Go project and CLI application skeleton is taken care
of using the `cobra` CLI utility program.

```bash
$ cobra init kubectl-ls --pkg-name ls --license MIT --author "Mario-Leander Reimer"
$ cd kubectl-ls
$ go mod init github.com/lreimer/go-for-operations/kubectl-ls

$ go install
$ kubectl ls
```

Open the generated `cmd/root.go` file, perform some cleanup and add the following code
to execute when the command is run.
```golang
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
```

## Adding a version subcommand

Adding new CLI commands is also performed using the `cobra` CLI utility program.

```bash
$ cobra add version
```

Open the generated `cmd/version.go` file, perform some cleanup and add the following code
to execute when the command is run.
```golang
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
	Short: "Display version",
	Long:  "Display version number and commit hash of application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kubectl-ls %v %v\n", version, commit)
	},
}
``` 

## Building and Releasing

As usual a `Makefile` with the basic build goals may be added to ease development. Building the
final distributions artefacts can easily be performed using GoReleaser.
