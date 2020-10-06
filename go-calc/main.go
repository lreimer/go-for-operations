package main

import (
	"flag"
	"fmt"
	"os"
)

var version = ""
var commit = ""

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addEnabled := addCmd.Bool("enabled", true, "enabled")

	if len(os.Args) < 2 {
		illegalArguments()
	}

	switch os.Args[1] {
	case "add":
		// parse the remaing arguments
		addCmd.Parse(os.Args[2:])
		if *addEnabled {
			args := addCmd.Args()
			result := Add(args[0], args[1])
			fmt.Printf("%v + %v = %v\n", args[0], args[1], result)
		}
	case "version":
		fmt.Printf("go-calc %v %v\n", version, commit)
	default:
		illegalArguments()
	}
}

func illegalArguments() {
	fmt.Println("Expected 'add' or 'version' subcommands.")
	os.Exit(1)
}
