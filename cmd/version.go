package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of devsync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devsync v%s\n", Version)
	},
}
