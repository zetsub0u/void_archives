package cmd

import (
	"fmt"

	"runtime"

	"github.com/spf13/cobra"
)

// These variables are populated via the Go linker.
var (
	version = "unknown"
	commit  = "unknown"
	branch  = "unknown"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("void archives %s %s %s %s (commit %s, branch %s)\n",
			version, runtime.GOOS, runtime.GOARCH, runtime.Version(), commit, branch)
	},
}
