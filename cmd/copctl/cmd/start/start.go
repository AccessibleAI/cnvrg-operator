package start

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Short:   "Start copctl services",
	}
)
