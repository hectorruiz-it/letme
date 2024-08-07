package config

import (
	"fmt"
	"os"

	utils "github.com/lockedinspace/letme/pkg"
	"github.com/spf13/cobra"
)

var GetContexts = &cobra.Command{
	Use: "get-contexts",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.LetmeConfigCreate()
		utils.ConfigFileHealth()
	},
	Short: "Get active and available contexts.",
	Long:  `List all configured contexts in your letme-config file marking the active context with '*'`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		contexts := utils.GetAvalaibleContexts()
		fmt.Println("Active context marked with '*': ")
		for _, context := range contexts {
			currentContext := utils.GetCurrentContext()
			if context == currentContext {
				fmt.Println("* " + context)
			} else {
				fmt.Println("  " + context)
			}
		}
		os.Exit(0)
	},
}

func init() {
	ConfigCmd.AddCommand(GetContexts)
}
