package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "copctl",
		Short: "cnvrg operator ctl",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	initZapLog()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COPCTL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initZapLog() {
	config := zap.NewDevelopmentConfig()
	//config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()

	zap.ReplaceGlobals(logger)

}
