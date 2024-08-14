package main

import (
	"fmt"
	v1alpha1 "github.com/AccessibleAI/cnvrg-operator/api/v1alpha1"
	"github.com/AccessibleAI/cnvrg-operator/controllers/metastorageprovider"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"strings"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	startOperatorCmd.PersistentFlags().StringP(
		"metrics-addr", "", ":8080", "The address the metric endpoint binds to.")
	startOperatorCmd.PersistentFlags().StringP(
		"health-probe-addr", "", ":8081",
		"The address the health probes endpoints (/healthz, /readyz) binds to.")
	startOperatorCmd.PersistentFlags().BoolP(
		"enable-leader-election", "", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	startOperatorCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	startOperatorCmd.PersistentFlags().IntP(
		"max-concurrent-reconciles", "", 1, "Max concurrent reconciles")
	startOperatorCmd.PersistentFlags().StringP(
		"namespace", "", "cnvrg", "the ns into which the operator has been deployed")

	viper.BindPFlag("metrics-addr", startOperatorCmd.PersistentFlags().Lookup("metrics-addr"))
	viper.BindPFlag("health-probe-addr", startOperatorCmd.PersistentFlags().Lookup("health-probe-addr"))
	viper.BindPFlag("enable-leader-election", startOperatorCmd.PersistentFlags().Lookup("enable-leader-election"))
	viper.BindPFlag("verbose", startOperatorCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("max-concurrent-reconciles", startOperatorCmd.PersistentFlags().Lookup("max-concurrent-reconciles"))
	viper.BindPFlag("namespace", startOperatorCmd.PersistentFlags().Lookup("namespace"))
}

// Root cmd and params
var rootCmd = &cobra.Command{
	Use:   "cnvrg-meta-storage-provider",
	Short: "cnvrg-meta-storage-provider - K8s operator for deploying storage for cnvrg stack",
}

var startOperatorCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cnvrg meta storage provider",
	Run: func(cmd *cobra.Command, args []string) {
		runOperator()
	},
}

func runOperator() {
	ctrl.SetLogger(zapr.NewLogger(initZapLog()))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{
				viper.GetString("namespace"): {},
			},
		},
		Metrics: metricsserver.Options{
			BindAddress: viper.GetString("metrics-addr"),
		},
		HealthProbeBindAddress: viper.GetString("health-probe-addr"),
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&metastorageprovider.Reconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		Namespace: viper.GetString("namespace"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CnvrgMetaStorageProvider")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupCommands() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CNVRG_META_STORAGE_PROVIDER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	rootCmd.AddCommand(startOperatorCmd)
}

func main() {
	setupCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if viper.GetBool("verbose") {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	logger, _ := config.Build()
	return logger
}
