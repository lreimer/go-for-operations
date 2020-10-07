package main

import "github.com/lreimer/go-for-operations/kubectl-chaos/cmd"

var version string
var commit string

func main() {
	cmd.Execute()
}
