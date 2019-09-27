package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"oci-tools/utils"
	"strconv"
)

var consoleCmd = &cobra.Command{
	Use: "console",
	Short: "instance console information",
	Long: `Get an instances console data. Instead of the two step process of the standard
CLI it takes care of it all in a single 'display' verb.
`,
}

var dispConsoleCmd = &cobra.Command{
	Use:   "display",
	Short: "display an instances console output",
	Long: `capture and display an instances console:

oci-tool console display --instance-ocid {ocid}

`,
	Run: func(cmd *cobra.Command, args []string) {
		displayConsole(cmd)
	},
}

func init() {
	computeCmd.AddCommand(consoleCmd)
	consoleCmd.AddCommand(dispConsoleCmd)
	dispConsoleCmd.Flags().String("instance-id", "", "instance id to get console output for")
	dispConsoleCmd.Flags().Int("length", 1000000, "lenght of console to fetch")
	_ = dispConsoleCmd.MarkFlagRequired("instance-id")
}

func displayConsole(cmd *cobra.Command) {
	ch, err := config.CaptureConsole(cmd.Flag("instance-id").Value.String())
	if err != nil {
		utils.ErrorMsg("error capturing console history", err)
	} else {
		l, err := strconv.Atoi(cmd.Flag("length").Value.String())
		if err != nil {
			utils.ErrorMsg("invalid length specified", err)
			return
		}
		cho, err := config.CollectConsole(ch, l)
		if err != nil {
			utils.ErrorMsg("error collecting console content", err)
		} else {
			fmt.Println(cho)
		}
	}
}


// 	length := 1000000