package create

import "github.com/spf13/cobra"

var (
	Cmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Helper for generating different assets",
	}
)
