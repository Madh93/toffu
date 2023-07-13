package cmd

import (
	"log"

	"github.com/Madh93/toffu/internal/toffu"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Show current status",
	Run: func(cmd *cobra.Command, args []string) {
		toffu := toffu.New()
		if err := toffu.GetStatus(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
