package cmd

import (
	"fmt"

	"github.com/Piya-Boy/devsync/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of devsync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devsync v%s\n", version.Version)
	},
}
