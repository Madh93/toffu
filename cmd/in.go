package cmd

import (
	"log"

	"github.com/Madh93/toffu/internal/toffu"
	"github.com/spf13/cobra"
)

var inCmd = &cobra.Command{
	Use:     "in",
	Aliases: []string{"i"},
	Short:   "On-demand clock-in",
	Run: func(cmd *cobra.Command, args []string) {
		toffu := toffu.New()
		if err := toffu.ClockIn(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(inCmd)
}
