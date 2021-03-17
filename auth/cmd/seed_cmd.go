package cmd

import (
	"github.com/spf13/cobra"
)

var seedCMD = cobra.Command{
	Use:  "seed database",
	Long: "seed database strucutures. This will seed tables",
	Run:  seed,
}

// seed database with fake data
func seed(cmd *cobra.Command, args []string) {}
