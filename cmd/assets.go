package cmd

import (
	"github.com/spf13/cobra"
)

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Manage assets within the masonictempl project.",
	Long: `Manage assets within the masonictempl project. For example:
You can use this command to manage remote assets.
masonictempl assets [command]`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(assetsCmd)
	assetsCmd.AddCommand(migrateToRemoteCmd)
}
