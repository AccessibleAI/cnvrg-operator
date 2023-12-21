package start

import (
	"crypto/tls"
	"fmt"
	"github.com/AccessibleAI/cnvrg-operator/pkg/admission"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

type PlatformType string

const (
	CAGPPlatform PlatformType = "cagp"
)

func init() {
	admissionCtrlCmd.PersistentFlags().StringP("platform", "", "cagp",
		fmt.Sprintf("one of: %s|", CAGPPlatform))
	admissionCtrlCmd.PersistentFlags().StringP("domain-pool", "", "zerossl", "domain pool to use when deploying on CAGP")
	admissionCtrlCmd.PersistentFlags().StringP("domain-claim", "", "cnvrg-domain-claim", "domain claim to create")
	admissionCtrlCmd.PersistentFlags().StringP("gateway-name", "", "cnvrg-gateway", "gateway to create")
	admissionCtrlCmd.PersistentFlags().StringP("namespace", "", "cnvrg", "namespace")
	admissionCtrlCmd.PersistentFlags().StringP("reg-user", "", "reg-user", "registry user")
	admissionCtrlCmd.PersistentFlags().StringP("reg-password", "", "reg-password", "registry password")
	admissionCtrlCmd.PersistentFlags().StringP("crt", "", "", "path to certificate file")
	admissionCtrlCmd.PersistentFlags().StringP("key", "", "", "path to key file")

	viper.BindPFlag("platform", admissionCtrlCmd.PersistentFlags().Lookup("platform"))
	viper.BindPFlag("domain-pool", admissionCtrlCmd.PersistentFlags().Lookup("domain-pool"))
	viper.BindPFlag("domain-claim", admissionCtrlCmd.PersistentFlags().Lookup("domain-claim"))
	viper.BindPFlag("namespace", admissionCtrlCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("reg-user", admissionCtrlCmd.PersistentFlags().Lookup("reg-user"))
	viper.BindPFlag("reg-password", admissionCtrlCmd.PersistentFlags().Lookup("reg-password"))
	viper.BindPFlag("crt", admissionCtrlCmd.PersistentFlags().Lookup("crt"))
	viper.BindPFlag("key", admissionCtrlCmd.PersistentFlags().Lookup("key"))

	Cmd.AddCommand(admissionCtrlCmd)
}

var admissionCtrlCmd = &cobra.Command{
	Use:     "admission-controller",
	Aliases: []string{"a"},
	Short:   "start cnvrg's operator K8s dynamic admission controller",
	Run: func(cmd *cobra.Command, args []string) {
		cert := viper.GetString("crt")
		key := viper.GetString("key")
		pair, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			zap.S().Error("Failed to load key pair: %v", err)
		}

		// Handler for CnvrgCap clusterDomain deployed on AI Cloud
		//http.HandleFunc("/cap/clusterdomain/mutate", admission.MutateCnvrgAppClusterDomainHandler)
		aiCloudDomainDiscoveryHandler := admission.NewAICloudDomainHandler()
		http.HandleFunc(aiCloudDomainDiscoveryHandler.HandlerPath(), aiCloudDomainDiscoveryHandler.Handler)

		// Create HTTPS server configuration
		s := &http.Server{
			Addr:      "0.0.0.0:8080",
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		}
		zap.S().Info("Admission controller started")
		zap.S().Fatal(s.ListenAndServeTLS("", ""))
	},
}
