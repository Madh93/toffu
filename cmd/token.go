package cmd

import (
	"log"

	"github.com/Madh93/toffu/internal/toffu"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"t"},
	Short:   "Generate API token (Required the first time)",
	Run: func(cmd *cobra.Command, args []string) {
		toffu := toffu.New()
		if err := toffu.GenerateToken(); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
