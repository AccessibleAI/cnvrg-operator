package main

import (
	"github.com/AccessibleAI/cnvrg-operator/pkg/dumper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type dump struct {
	dumper.Dumper
}

var (
	profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "profile - list and dump cnvrg deployment profiles",
	}

	profileDumpParams = []param{
		{name: "templates-dump-dir", shorthand: "d", value: "cnvrg-manifests", usage: "dump cnvrg stack components"},
		{name: "preserve-templates-dir", shorthand: "k", value: false, usage: "preserve templates directories"},
		{name: "namespace", shorthand: "n", value: "cnvrg", usage: "set K8s namespace"},
	}
	profileDumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "dump - dump cnvrg deployment manifests",
	}

	dumpControlPlaneParams = []param{
		{name: "ingress", value: "ingress", shorthand: "s", usage: "must be one of the: istio|ingress|openshift|nodeport"},
		{name: "https", value: false, usage: "enabled|disable https"},
		{name: "proxy", value: false, usage: "enable proxy"},
		{name: "wildcard-domain", shorthand: "a", value: "", usage: "the wildcard domain for cnvrg stack deployments"},
		{name: "control-plane-image", shorthand: "i", value: "", usage: "cnvrg control plane image"},
		{name: "registry-user", shorthand: "u", value: "", usage: "docker registry user"},
		{name: "registry-password", shorthand: "p", value: "", usage: "docker registry password"},
		{name: "cri", shorthand: "c", value: "containerd", usage: "container runtime interface one of: docker|containerd|cri-o"},
	}
	dumpControlPlaneCmd = &cobra.Command{
		Use:     "control-plane",
		Short:   "dump cnvrg control plane as raw K8s manifests",
		Aliases: []string{"cp"},
		Run: func(cmd *cobra.Command, args []string) {
			cp := dumper.NewControlPlane(
				viper.GetString("control-plane-image"),
				viper.GetString("wildcard-domain"),
				viper.GetString("cri"),
				viper.GetString("registry-user"),
				viper.GetString("registry-password"),
				viper.GetString("ingress"),
				viper.GetString("namespace"),
				viper.GetBool("https"),
				viper.GetBool("proxy"),
			)
			d := dump{cp}
			if err := d.BuildState(); err != nil {
				log.Fatal(err)
			}
			if err := d.Dump(viper.GetBool("preserve-templates-dir")); err != nil {
				log.Fatal(err)
			}
		},
	}

	dumpNetworkParams = []param{}
	dumpNetworkCmd    = &cobra.Command{
		Use:   "network",
		Short: "dump cnvrg control plane as raw K8s manifests",
		Run: func(cmd *cobra.Command, args []string) {
			network := dumper.NewNetworking(viper.GetString("namespace"))
			d := dump{network}
			if err := d.BuildState(); err != nil {
				log.Fatal(err)
			}
			if err := d.Dump(viper.GetBool("preserve-templates-dir")); err != nil {
				log.Fatal(err)
			}
		},
	}
)
