package main

import "github.com/lreimer/go-for-operations/kubectl-ls/cmd"

var version = ""

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
