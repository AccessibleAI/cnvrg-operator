package main

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/controllers"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	scheme            = runtime.NewScheme()
	setupLog          = ctrl.Log.WithName("setup")
	runOperatorParams = []param{
		{name: "metrics-addr", shorthand: "", value: ":8080", usage: "The address the metric endpoint binds to."},
		{name: "enable-leader-election", shorthand: "", value: false, usage: "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager."},
		{name: "dry-run", shorthand: "", value: false, usage: "Only parse templates, without applying"},
		{name: "deploy-depended-crds", shorthand: "", value: true, usage: "Deploy depended (external) CRDs automatically"},
		{name: "own-istio-resources", shorthand: "", value: true, usage: "Watch for istio resources"},
		{name: "own-openshift-resources", shorthand: "", value: false, usage: "Watch for OpenShift resources"},
		{name: "own-prometheus-resources", shorthand: "", value: true, usage: "Watch for Prometheus resources"},
	}
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = mlopsv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func initZapLog() *zap.Logger {

	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	//config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}

// run operator cmd and params
var runOperatorCmd = &cobra.Command{
	Use:   "run",
	Short: "Run cnvrg operator",
	Run: func(cmd *cobra.Command, args []string) {
		loggerMgr := initZapLog()
		zap.ReplaceGlobals(loggerMgr)
		logger := loggerMgr.Sugar()
		logger.Info("this is info")
		logger.Debug("this is debug")
		ensureCrdsAvailability()
		os.Exit(1)

		runOperator()
	},
}

func ensureCrdsAvailability() {
	if viper.GetBool("own-istio-resources") {
		networking.LoadCrds()
	}
}

func runOperator() {
	//zapLog, err := zap.NewDevelopment()
	//ctrl.SetLogger(zapr.NewLogger(zapLog))
	l := initZapLog()
	l.Sugar()
	ctrl.SetLogger(zapr.NewLogger(l))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: viper.GetString("metrics-addr"),
		Port:               9443,
		LeaderElection:     viper.GetBool("enable-leader-election"),
		LeaderElectionID:   "99748453.cnvrg.io",
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

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
