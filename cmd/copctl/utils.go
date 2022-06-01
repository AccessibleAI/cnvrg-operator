package main

import (
	"fmt"
	"github.com/AccessibleAI/cnvrg-operator/pkg/dumper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type PlainFormatter struct{}

func setParams(params []*dumper.Param, command *cobra.Command) {
	for _, param := range params {
		switch v := param.Value.(type) {
		case int:
			command.PersistentFlags().IntP(param.Name, param.Shorthand, v, param.Usage)
		case string:
			command.PersistentFlags().StringP(param.Name, param.Shorthand, v, param.Usage)
		case bool:
			command.PersistentFlags().BoolP(param.Name, param.Shorthand, v, param.Usage)
		}
		if err := viper.BindPFlag(param.Name, command.PersistentFlags().Lookup(param.Name)); err != nil {
			panic(err)
		}
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CNVRGCTL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setupLogging()

}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func setupLogging() {

	log.SetFormatter(new(PlainFormatter))
	log.SetOutput(os.Stdout)
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}
}
