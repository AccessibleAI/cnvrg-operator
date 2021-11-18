package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"strings"

	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
)

type param struct {
	name      string
	shorthand string
	value     interface{}
	usage     string
	required  bool
}

var (
	BuildVersion string
)

var (
	scheme            = runtime.NewScheme()
	setupLog          = ctrl.Log.WithName("setup")
	rootParams        = []param{}
	runOperatorParams = []param{
		{name: "metrics-addr", shorthand: "", value: ":8080", usage: "The address the metric endpoint binds to."},
		{name: "health-probe-addr", shorthand: "", value: ":8081", usage: "The address the health probes endpoints (/healthz, /readyz) binds to."},
		{name: "enable-leader-election", shorthand: "", value: false, usage: "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager."},
		{name: "dry-run", shorthand: "", value: false, usage: "Only parse templates, without applying"},
		{name: "templates-dump-dir", shorthand: "", value: "", usage: "destination dir for rendering templates for debugging"},
		{name: "verbose", shorthand: "v", value: false, usage: "Verbose output"},
		{name: "deploy-depended-crds", shorthand: "", value: true, usage: "Deploy depended (external) CRDs automatically"},
		{name: "own-istio-resources", shorthand: "", value: true, usage: "Watch for istio resources"},
		{name: "own-openshift-resources", shorthand: "", value: false, usage: "Watch for OpenShift resources"},
		{name: "own-prometheus-resources", shorthand: "", value: true, usage: "Watch for Prometheus resources"},
		{name: "max-concurrent-reconciles", shorthand: "", value: 1, usage: "Max concurrent reconciles"},
		{name: "cleanup-pvc", shorthand: "", value: false, usage: "set to true to delete PVCs on CR delete"},
	}
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = mlopsv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

// Root cmd and params
var rootCmd = &cobra.Command{
	Use:   "cnvrg-operator",
	Short: "cnvrg-operator - K8s operator for deploying cnvrg stack",
}

var runOperatorCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cnvrg operator",
	Run: func(cmd *cobra.Command, args []string) {
		loggerMgr := initZapLog()
		loggerMgr.Sugar()
		zap.ReplaceGlobals(loggerMgr)
		runOperator()
	},
}

var operatorVersion = &cobra.Command{
	Use:   "version",
	Short: "Print cnvrg operator version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🐾 version: %s\n", BuildVersion)
	},
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

func setParams(params []param, command *cobra.Command) {
	for _, param := range params {
		switch v := param.value.(type) {
		case int:
			command.PersistentFlags().IntP(param.name, param.shorthand, v, param.usage)
		case string:
			command.PersistentFlags().StringP(param.name, param.shorthand, v, param.usage)
		case bool:
			command.PersistentFlags().BoolP(param.name, param.shorthand, v, param.usage)
		}
		if err := viper.BindPFlag(param.name, command.PersistentFlags().Lookup(param.name)); err != nil {
			panic(err)
		}
	}
}

func runOperator() {
	ctrl.SetLogger(zapr.NewLogger(initZapLog()))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     viper.GetString("metrics-addr"),
		HealthProbeBindAddress: viper.GetString("health-probe-addr"),
		Port:                   9443,
		LeaderElection:         viper.GetBool("enable-leader-election"),
		LeaderElectionID:       "99748453.cnvrg.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.CnvrgAppReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("CnvrgApp"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CnvrgApp")
		os.Exit(1)
	}

	if err = (&controllers.CnvrgInfraReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("CnvrgInfra"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CnvrgInfra")
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

	zap.S().Infof("cnvrg operator version: %s", BuildVersion)
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupCommands() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CNVRG_OPERATOR")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Setup commands
	setParams(runOperatorParams, runOperatorCmd)
	setParams(rootParams, rootCmd)
	rootCmd.AddCommand(operatorVersion)
	rootCmd.AddCommand(runOperatorCmd)
}

func informPkger() {
	pkger.Include("/pkg/controlplane/tmpl")
	pkger.Include("/pkg/dbs/tmpl")
	pkger.Include("/pkg/logging/tmpl")
	pkger.Include("/pkg/monitoring/tmpl")
	pkger.Include("/pkg/networking/tmpl")
	pkger.Include("/pkg/registry/tmpl")
	pkger.Include("/pkg/storage/tmpl")
	pkger.Include("/pkg/gpu/tmpl")
	pkger.Include("/pkg/reloader/tmpl")
	pkger.Include("/pkg/capsule/tmpl")
	pkger.Include("/pkg/priorityclass/tmpl")
	pkger.Include("/pkg/ingresscheck/tmpl")
}

func main() {
	informPkger()
	setupCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
