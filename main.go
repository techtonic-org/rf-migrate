package main

import (
	"fmt"
	"os"

	"github.com/techtonic-org/rf-migrate/cmd"
)

// version is set during build using ldflags
var version = "dev"

func main() {
	// Check for version flag
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("rf-migrate version %s\n", version)
			os.Exit(0)
		}
	}

	cmd.Execute()
}
