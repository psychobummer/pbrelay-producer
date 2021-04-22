package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pbrelay-producer",
	Short: "A tool for sending arbitrary data to a PsychoBummer(t)(r)(tm) relayserver",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func init() {
}
