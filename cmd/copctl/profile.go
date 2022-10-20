package main

import (
	"encoding/json"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers/app"
	"github.com/AccessibleAI/cnvrg-operator/pkg/dumper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			a := mlopsv1.DefaultCnvrgAppSpec()
			spec := &mlopsv1.CnvrgApp{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-app", Namespace: "cnvrg"}}

			if err := app.CalculateAndApplyAppDefaults(spec, &a, nil); err != nil {
				log.Fatal(err)
			}

			spec.Spec = a

			b, err := json.Marshal(spec)
			if err != nil {
				log.Fatal(err)
			}
			log.Info(string(b))
		},
	}

	dumpNetworkParams = []param{}
	dumpNetworkCmd    = &cobra.Command{
		Use:   "network",
		Short: "dump cnvrg control plane as raw K8s manifests",
		Run: func(cmd *cobra.Command, args []string) {
			//network := dumper.NewNetworking(viper.GetString("namespace"))
			//d := dump{network}
			//if err := d.BuildState(); err != nil {
			//	log.Fatal(err)
			//}
			//if err := d.Dump(viper.GetBool("preserve-templates-dir")); err != nil {
			//	log.Fatal(err)
			//}
		},
	}
)
