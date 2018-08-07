package fecmd

import (
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	"aether-core/services/logging"
	"strings"
)

type flag struct {
	value   interface{}
	changed bool
}

type flags struct {
	loggingLevel flag // int
	clientIp     flag // string
	clientPort   flag // int

	// add more flags here
}

func renderFlags(cmd *cobra.Command) flags {
	var fl flags

	flg, err := cmd.Flags().GetInt("logginglevel")
	if err != nil && !strings.Contains(err.Error(), "flag accessed but not defined") {
		logging.LogCrash(err)
	}
	fl.loggingLevel.value = flg
	fl.loggingLevel.changed = cmd.Flags().Changed("logginglevel")

	flg2, err2 := cmd.Flags().GetString("clientip")
	if err2 != nil && !strings.Contains(err2.Error(), "flag accessed but not defined") {
		logging.LogCrash(err2)
	}
	fl.clientIp.value = flg2
	fl.clientIp.changed = cmd.Flags().Changed("clientip")

	flg3, err3 := cmd.Flags().GetInt("clientport")
	if err3 != nil && !strings.Contains(err3.Error(), "flag accessed but not defined") {
		logging.LogCrash(err3)
	}
	fl.clientPort.value = flg3
	fl.clientPort.changed = cmd.Flags().Changed("clientport")

	// add more flags here

	return fl
}
