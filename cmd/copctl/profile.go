package main

import (
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/controlplane"
	"github.com/AccessibleAI/cnvrg-operator/pkg/dbs"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/monitoring"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/priorityclass"
	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
	"github.com/Dimss/crypt/apr1_crypt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type Dumper interface {
	BuildState() []*desired.State
}

var (
	profileDumpParams = []param{
		{name: "ingress", value: "ingress", usage: "must be one of the: istio|ingress|openshift|nodeport"},
		{name: "wildcard-domain", shorthand: "w", value: "", usage: "the wildcard domain for cnvrg stack deployments"},
		{name: "control-plane", value: false, usage: "dump cnvrg control plane stack"},
		{name: "control-plane-image", shorthand: "i", value: "", usage: "cnvrg control plane image"},
		{name: "cri", shorthand: "c", value: "containerd", usage: "container runtime interface one of: docker|containerd|cri-o"},
		{name: "registry-user", shorthand: "u", value: "", usage: "docker registry user"},
		{name: "registry-password", shorthand: "p", value: "", usage: "docker registry password"},
		{name: "networking", value: false, usage: "dump cnvrg networking stack"},
		{name: "logging", value: false, usage: "dump cnvrg logging stack"},
		{name: "monitoring", value: false, usage: "dump cnvrg monitoring stack"},
		{name: "templates-dump-dir", shorthand: "d", value: "./cnvrg-manifests", usage: "dump cnvrg monitoring stack"},
	}
	profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "profile - list and dump cnvrg deployment profiles",
	}

	profileList = &cobra.Command{
		Use:   "list",
		Short: "list - list available deployment profiles",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Cnvrg configuration profiles:")
		},
	}

	profileDump = &cobra.Command{
		Use:   "dump",
		Short: "dump - dump cnvrg deployment manifests",
	}

	dumpControlPlane = &cobra.Command{
		Use:   "control-plane",
		Short: "dump cnvrg control plane as raw K8s manifests",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func dumpCnvrgNetworkingStack() {
	infra := &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}
	infra.Spec.InfraNamespace = "istio-test"
	if err := controllers.CalculateAndApplyInfraDefaults(infra, &infra.Spec, nil); err != nil {
		log.Fatal(err)
	}

	infra.Spec.Networking.Istio.Enabled = true
	infra.Spec.CnvrgAppPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-apps", Value: 2000000, Description: "cnvrg control plane apps priority class"}
	state := networking.InfraNetworkingState(infra)
	state = append(state, networking.IstioCrds()...)
	state = append(state, priorityclass.State()...)

	if err := os.RemoveAll(viper.GetString("templates-dump-dir")); err != nil {
		log.Fatal(err)
	}

	for _, o := range state {

		if o.TemplateData == nil {
			o.TemplateData = infra
		}
		if err := o.GenerateDeployable(); err != nil {
			log.Fatal(err)
		}
		if err := o.DumpTemplateToFile(); err != nil {
			log.Fatal(err)
		}
	}

}

func dumpCnvrgControlPlane() {
	app := &mlopsv1.CnvrgApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cnvrg-app",
			Namespace: "cnvrg",
		},
		Spec: mlopsv1.DefaultCnvrgAppSpec(),
	}
	infra := &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}
	if err := controllers.CalculateAndApplyAppDefaults(app, &app.Spec, infra, nil); err != nil {
		log.Fatal(err)
	}
	app.Spec.Annotations = make(map[string]string)
	app.Spec.Annotations["purpose"] = "cnvrg"
	app.Spec.ControlPlane.WebApp.Enabled = true

	app.Spec.ControlPlane.BaseConfig.FeatureFlags = make(map[string]string)
	app.Spec.ControlPlane.BaseConfig.FeatureFlags["CNVRG_MOUNT_HOST_FOLDERS"] = "false"

	app.Spec.ControlPlane.Image = viper.GetString("control-plane-image")
	app.Spec.Cri = mlopsv1.CriType(viper.GetString("cri"))

	app.Spec.Networking.Ingress.Type = mlopsv1.IngressType(viper.GetString("ingress"))
	if app.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress {

		app.Spec.Networking.Ingress.IstioGwEnabled = true
	}
	app.Spec.ClusterDomain = viper.GetString("wildcard-domain")

	app.Spec.ControlPlane.Sidekiq.Enabled = true
	app.Spec.ControlPlane.Sidekiq.Split = true
	app.Spec.ControlPlane.Hyper.Enabled = true
	app.Spec.ControlPlane.CnvrgScheduler.Enabled = false
	app.Spec.ControlPlane.Searchkiq.Enabled = true
	app.Spec.ControlPlane.Systemkiq.Enabled = true

	app.Spec.ControlPlane.ObjectStorage.Type = "minio"
	app.Spec.ControlPlane.ObjectStorage.AzureContainer = "azure"
	app.Spec.ControlPlane.ObjectStorage.AzureAccountName = "azure"
	app.Spec.ControlPlane.ObjectStorage.AzureContainer = "azure"
	app.Spec.ControlPlane.ObjectStorage.AccessKey = "AKIAIOSFODNN7EXAMPLE"
	app.Spec.ControlPlane.ObjectStorage.SecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	app.Spec.ControlPlane.ObjectStorage.GcpProject = "gcp-project"

	app.Spec.ControlPlane.SMTP.Port = 25
	app.Spec.ControlPlane.SMTP.Server = "smtp"
	app.Spec.ControlPlane.SMTP.Username = "username"
	app.Spec.ControlPlane.SMTP.Password = "password"
	app.Spec.ControlPlane.SMTP.Domain = "domain"
	app.Spec.ControlPlane.SMTP.OpensslVerifyMode = "false"

	app.Spec.ControlPlane.Ldap.Host = "host"
	app.Spec.ControlPlane.Ldap.Port = "port"
	app.Spec.ControlPlane.Ldap.Ssl = "false"
	app.Spec.ControlPlane.Ldap.Base = "base"
	app.Spec.ControlPlane.Ldap.AdminUser = "admin"
	app.Spec.ControlPlane.Ldap.AdminPassword = "pass"

	app.Spec.ControlPlane.BaseConfig.SentryURL = "sentry-url"

	app.Spec.Dbs.Es.Enabled = true
	app.Spec.Dbs.Pg.Enabled = true
	app.Spec.Dbs.Minio.Enabled = true
	app.Spec.Dbs.Redis.Enabled = true
	app.Spec.CnvrgAppPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-apps", Value: 2000000, Description: "cnvrg control plane apps priority class"}
	app.Spec.CnvrgJobPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-jobs", Value: 1000000, Description: "cnvrg jobs priority class"}

	redisSecretData := desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":   app.Namespace,
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
			"CredsRef":    app.Spec.Dbs.Redis.CredsRef,
			"SvcName":     app.Spec.Dbs.Redis.SvcName,
		},
	}

	pgSecretData := desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":          app.Namespace,
			"CredsRef":           app.Spec.Dbs.Pg.CredsRef,
			"Annotations":        app.Spec.Annotations,
			"Labels":             app.Spec.Labels,
			"MaxConnections":     app.Spec.Dbs.Pg.MaxConnections,
			"SharedBuffers":      app.Spec.Dbs.Pg.SharedBuffers,
			"EffectiveCacheSize": app.Spec.Dbs.Pg.EffectiveCacheSize,
			"SvcName":            app.Spec.Dbs.Pg.SvcName,
		},
	}

	esSecretData := desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":   app.Namespace,
			"CredsRef":    app.Spec.Dbs.Es.CredsRef,
			"EsUrl":       fmt.Sprintf("%s.%s.svc:%d", app.Spec.Dbs.Es.SvcName, app.Namespace, app.Spec.Dbs.Es.Port),
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
		},
	}

	user := "cnvrg"
	pass := desired.RandomString()
	passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
	if err != nil {
		log.Fatal(err)
	}

	promSecretData := desired.TemplateData{
		Data: map[string]interface{}{
			"Namespace":   app.Namespace,
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
			"CredsRef":    app.Spec.Monitoring.Prometheus.CredsRef,
			"User":        user,
			"Pass":        pass,
			"PassHash":    fmt.Sprintf("%s:%s", user, passHash),
			"PromUrl":     fmt.Sprintf("http://%s.%s.svc:%d", app.Spec.Monitoring.Prometheus.SvcName, app.Namespace, app.Spec.Monitoring.Prometheus.Port),
		},
	}

	app.Spec.Registry.User = viper.GetString("registry-user")
	app.Spec.Registry.Password = viper.GetString("registry-password")

	registryData := desired.TemplateData{
		Namespace: app.Namespace,
		Data: map[string]interface{}{
			"Registry":    app.Spec.Registry,
			"Annotations": app.Spec.Annotations,
			"Labels":      app.Spec.Labels,
		},
	}

	state := controlplane.State(app)
	state = append(state, dbs.AppDbsState(app)...)
	state = append(state, networking.CnvrgAppNetworkingState(app)...)
	state = append(state, priorityclass.State()...)
	state = append(state, dbs.RedisCreds(redisSecretData)...)
	state = append(state, dbs.PgCreds(pgSecretData)...)
	state = append(state, dbs.EsCreds(esSecretData)...)
	state = append(state, monitoring.PromCreds(promSecretData)...)
	state = append(state, registry.State(registryData)...)

	if err := os.RemoveAll(viper.GetString("templates-dump-dir")); err != nil {
		log.Fatal(err)
	}

	for _, o := range state {

		if o.TemplateData == nil {
			o.TemplateData = app
		}
		if err := o.GenerateDeployable(); err != nil {
			log.Fatal(err)
		}
		if err := o.DumpTemplateToFile(); err != nil {
			log.Fatal(err)
		}
	}
}
