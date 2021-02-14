package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type param struct {
	name      string
	shorthand string
	value     interface{}
	usage     string
	required  bool
}

// Root cmd and params
var rootCmd = &cobra.Command{
	Use:   "cnvrg-operator",
	Short: "cnvrg-operator - K8s operator for deploying cnvrg stack",
}
var rootParams = []param{}

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

func setupCommands() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Setup commands
	setParams(runOperatorParams, runOperatorCmd)
	setParams(rootParams, rootCmd)
	rootCmd.AddCommand(runOperatorCmd)
}

func main() {
	setupCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
