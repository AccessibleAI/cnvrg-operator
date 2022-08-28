package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type PlainFormatter struct{}

type param struct {
	name      string
	shorthand string
	value     interface{}
	usage     string
	required  bool
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
