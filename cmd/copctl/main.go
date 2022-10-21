package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	Version    string
	Build      string
	rootParams = []param{}

	rootCmd = &cobra.Command{
		Use:   "copctl",
		Short: "copctl - cnvrg operator control tool",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	setParams(rootParams, rootCmd)
	setParams(profileDumpParams, profileDumpCmd)
	setParams(dumpControlPlaneParams, dumpControlPlaneCmd)
	setParams(dumpNetworkParams, dumpNetworkCmd)

	profileCmd.AddCommand(profileDumpCmd)
	profileDumpCmd.AddCommand(dumpControlPlaneCmd)
	profileDumpCmd.AddCommand(dumpNetworkCmd)
	rootCmd.AddCommand(profileCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("failed to start capsule-agent, %s", err)
		os.Exit(1)
	}
}
