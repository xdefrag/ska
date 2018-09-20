package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xdefrag/ska/internal/ska"
)

var rootCmd = &cobra.Command{
	Use:   "ska [template] [output]",
	Short: "Render template to output",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			template = templatesDir + "/" + args[0]
			output   = args[1]
		)

		_, err := os.Stat(template)
		handleError(err)

		vv, err := ska.ParseValues(template + "/values.toml")
		handleError(err)

		err = ska.GenerateTemplates(template+"/templates", output, vv)
		handleError(err)

		os.Exit(0)
	},
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}
