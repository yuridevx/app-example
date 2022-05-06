package cmd

import (
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "example",
	}

	cmd.AddCommand(
		run(),
	)

	return cmd
}

func Execute() {
	root := rootCmd()
	cobra.CheckErr(root.Execute())
}
