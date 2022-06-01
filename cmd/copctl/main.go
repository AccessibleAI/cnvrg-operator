package main

import (
	"github.com/AccessibleAI/cnvrg-operator/pkg/dumper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	Version    string
	Build      string
	rootParams = []*dumper.Param{}

	rootCmd = &cobra.Command{
		Use:   "copctl",
		Short: "copctl - cnvrg operator control tool",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	setParams(rootParams, rootCmd)
	setParams(dumper.NewControlPlane().GetCliFlags(), profileDump)
	profileCmd.AddCommand(profileList)
	profileCmd.AddCommand(profileDump)
	profileDump.AddCommand(dumpControlPlane)
	rootCmd.AddCommand(profileCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("failed to start capsule-agent, %s", err)
		os.Exit(1)
	}
}
