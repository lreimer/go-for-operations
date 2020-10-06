package main

import "github.com/lreimer/go-for-operations/kubectl-ls/cmd"

var version = ""
var commit = ""

func main() {
	cmd.SetVersion(version)
	cmd.SetCommit(commit)
	cmd.Execute()
}
