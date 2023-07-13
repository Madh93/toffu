package cmd

import (
	"log"

	"github.com/Madh93/toffu/internal/toffu"
	"github.com/spf13/cobra"
)

var outCmd = &cobra.Command{
	Use:     "out",
	Aliases: []string{"o"},
	Short:   "On-demand clock-out",
	Run: func(cmd *cobra.Command, args []string) {
		toffu := toffu.New()
		if err := toffu.ClockOut(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(outCmd)
}
