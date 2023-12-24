package get

import "github.com/spf13/cobra"

var (
	Cmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get copctl resources",
	}
)
