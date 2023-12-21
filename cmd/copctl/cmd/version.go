package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version = "v1alpha1"
	Build   string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and git sha256",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🐾 version: %s build: %s\n", Version, Build)
	},
}
