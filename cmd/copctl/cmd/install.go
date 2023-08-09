package cmd

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"os"
	"os/signal"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strings"
	"syscall"
	"time"
)

type PlatformType string

const (
	CAGPPlatform PlatformType = "cagp"
)

func init() {
	installCmd.PersistentFlags().StringP("platform", "", "cagp",
		fmt.Sprintf("one of: %s|", CAGPPlatform))
	installCmd.PersistentFlags().StringP("domain-pool", "", "zerossl", "domain pool to use when deploying on CAGP")
	installCmd.PersistentFlags().StringP("domain-claim", "", "cnvrg-domain-claim", "domain claim to create")
	installCmd.PersistentFlags().StringP("gateway-name", "", "cnvrg-gateway", "gateway to create")
	installCmd.PersistentFlags().StringP("namespace", "", "cnvrg", "namespace")
	installCmd.PersistentFlags().StringP("reg-user", "", "reg-user", "registry user")
	installCmd.PersistentFlags().StringP("reg-password", "", "reg-password", "registry password")

	viper.BindPFlag("platform", installCmd.PersistentFlags().Lookup("platform"))
	viper.BindPFlag("domain-pool", installCmd.PersistentFlags().Lookup("domain-pool"))
	viper.BindPFlag("domain-claim", installCmd.PersistentFlags().Lookup("domain-claim"))
	viper.BindPFlag("namespace", installCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("reg-user", installCmd.PersistentFlags().Lookup("reg-user"))
	viper.BindPFlag("reg-password", installCmd.PersistentFlags().Lookup("reg-password"))

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install cnvrg on different platforms",
	Run: func(cmd *cobra.Command, args []string) {
		if PlatformType(viper.GetString("platform")) == CAGPPlatform {
			zap.S().Infof("installing cnvrg on %s platform", CAGPPlatform)
			dc := getDomainClaim(viper.GetString("domain-claim"), viper.GetString("namespace"))
			// if domain claim not exists, create it
			if dc == nil {
				createDomainClaim(
					viper.GetString("domain-claim"),
					viper.GetString("namespace"),
					viper.GetString("domain-pool"),
				)
			}

			waitForDomainClaimReady(viper.GetString("domain-claim"), viper.GetString("namespace"))
			dcData := domainClaimData(dc)

			cnvrgAppSpec := cnvrgAppDeploySpec(
				strings.ReplaceAll(dcData["commonName"], "*.", ""),
				viper.GetString("reg-user"),
				viper.GetString("reg-password"),
				fmt.Sprintf("%s-%s", dcData["domainRefName"], viper.GetString("domain-claim")),
				viper.GetString("namespace"),
			)

			applyCnvrgSpec(cnvrgAppSpec)
			zap.S().Info("done")

			// handle interrupts
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			for {
				select {
				case s := <-sigCh:
					zap.S().Infof("signal: %s, shutting down", s)
					zap.S().Info("bye bye 👋")
					os.Exit(0)
				}
			}

		}
	},
}

func waitForDomainClaimReady(dcName, ns string) {
	for {
		dc := getDomainClaim(dcName, ns)
		if dc != nil {
			status := dc.Object["status"].(map[string]interface{})
			if len(status) > 0 {
				if val, ok := status["status"].(string); ok {
					if val == "READY" {
						zap.S().Infof("domain claim %s is ready yet", dcName)
						return
					}

				}

			}
		}
		zap.S().Infof("domain claim %s not ready yet", dcName)
		time.Sleep(5 * time.Second)
	}
}

func domainClaimData(dc *unstructured.Unstructured) map[string]string {

	return map[string]string{
		"domainRefName": dc.Object["spec"].(map[string]interface{})["domainRef"].(map[string]interface{})["name"].(string),
		"commonName":    dc.Object["spec"].(map[string]interface{})["domainRef"].(map[string]interface{})["commonName"].(string),
	}

}

func createDomainClaim(name, ns, domainPool string) {
	zap.S().Info("creating domain claim")
	if _, err := dynamicSet().Resource(domainClaimGVR()).
		Namespace(ns).
		Create(context.Background(), domainClaimSpec(name, domainPool), metav1.CreateOptions{}); err != nil {
		zap.Error(err)
	}

}

func getDomainClaim(dcName, dcNamespace string) *unstructured.Unstructured {
	u, err := dynamicSet().
		Resource(domainClaimGVR()).
		Namespace(dcNamespace).
		Get(context.Background(), dcName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		zap.S().Info("domain name not exists")
		return nil
	} else if err != nil {
		zap.S().Error(err)
		return nil
	}
	return u
}

func domainClaimGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "metacloud.cnvrg.io", Version: "v1alpha1", Resource: "domainclaims"}
}

func cnvrgAppGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "mlops.cnvrg.io", Version: "v1", Resource: "cnvrgapps"}
}

func domainClaimSpec(name, domainPool string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": domainClaimGVR().GroupVersion().String(),
			"kind":       "DomainClaim",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"secretName": name,
				"domainRef": map[string]interface{}{
					"domainPool": domainPool,
				},
			},
		},
	}
}

func dynamicSet() *dynamic.DynamicClient {
	rc, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal(err)
	}

	dynamicset, err := dynamic.NewForConfig(rc)
	if err != nil {
		zap.S().Fatal(err)
	}

	return dynamicset
}

func applyCnvrgSpec(cap *mlopsv1.CnvrgApp) {
	ucap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cap)
	if err != nil {
		zap.S().Error(err)
		return
	}

	if _, err := dynamicSet().
		Resource(cnvrgAppGVR()).
		Namespace(cap.Namespace).
		Create(
			context.Background(),
			&unstructured.Unstructured{Object: ucap},
			metav1.CreateOptions{},
		); err != nil {
		zap.S().Error(err)
	}
}

func cnvrgAppDeploySpec(clusterDomain, regUser, regPass, certSecret, ns string) *mlopsv1.CnvrgApp {
	cnvrgApp := &mlopsv1.CnvrgApp{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CnvrgApp",
			APIVersion: "mlops.cnvrg.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cnvrg-app",
			Namespace: ns,
		},
		Spec: mlopsv1.DefaultCnvrgAppSpec(),
	}
	cnvrgApp.Spec.ClusterDomain = clusterDomain
	cnvrgApp.Spec.ControlPlane.Image = "app:v4.7.52-DEV-15824-cnvrg-agnostic-infra-38"
	cnvrgApp.Spec.ControlPlane.CnvrgScheduler.Enabled = false
	cnvrgApp.Spec.ControlPlane.BaseConfig.FeatureFlags = map[string]string{
		"CNVRG_ENABLE_MOUNT_FOLDERS": "false",
		"CNVRG_MOUNT_HOST_FOLDERS":   "false",
		"CNVRG_PROMETHEUS_METRICS":   "true",
	}
	cnvrgApp.Spec.Registry.User = regUser
	cnvrgApp.Spec.Registry.Password = regPass
	cnvrgApp.Spec.Networking.Ingress.IstioIngressSelectorValue = "cluster-gateway"
	cnvrgApp.Spec.Networking.HTTPS.Enabled = true
	cnvrgApp.Spec.Networking.HTTPS.CertSecret = certSecret

	cnvrgApp.Spec.ControlPlane.WebApp.Enabled = true
	cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled = true
	cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled = true
	cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled = true
	cnvrgApp.Spec.ControlPlane.Hyper.Enabled = true
	cnvrgApp.Spec.ControlPlane.Nomex.Enabled = true
	cnvrgApp.Spec.Dbs.Pg.Enabled = true
	cnvrgApp.Spec.Dbs.Minio.Enabled = true
	cnvrgApp.Spec.Dbs.Es.Enabled = true
	cnvrgApp.Spec.Dbs.Redis.Enabled = true
	cnvrgApp.Spec.Dbs.Prom.Enabled = true
	cnvrgApp.Spec.Dbs.Prom.Grafana.Enabled = true

	return cnvrgApp
}
