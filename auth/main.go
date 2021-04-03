package main

import (
	"blogfa/auth/cmd"
	"fmt"
	"os"
)

// root execute command with cobra
func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}
