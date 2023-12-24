package start

import (
	"crypto/tls"
	"github.com/AccessibleAI/cnvrg-operator/pkg/admission"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

func init() {

	admissionCtrlCmd.PersistentFlags().StringP("crt", "", "", "path to certificate file")
	admissionCtrlCmd.PersistentFlags().StringP("key", "", "", "path to key file")

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
		aiCloudDomainDiscoveryHandler := admission.NewAICloudDomainHandler()
		readiness := admission.NewReadinessHandler()
		http.HandleFunc(aiCloudDomainDiscoveryHandler.HandlerPath(), aiCloudDomainDiscoveryHandler.Handler)
		http.HandleFunc(readiness.HandlerPath(), readiness.Handler)

		addr := "0.0.0.0:8080"

		// Create HTTPS server configuration
		s := &http.Server{
			Addr:      addr,
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		}
		zap.S().Infof("admission controller started on %s", addr)
		zap.S().Fatal(s.ListenAndServeTLS("", ""))
	},
}
