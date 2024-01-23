package cmd

import (
	"github.com/spf13/cobra"
)

var blogsCmd = &cobra.Command{
	Use:   "blog",
	Short: "Manage blogs within the masonictempl project.",
	Long: `Manage blogs within the masonictempl project. For example:
You canuse this command to list, create, update, and delete blogs.
The default will be to list all blogs.
masonictempl blogs [command]`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(blogsCmd)
	blogsCmd.AddCommand(blogCreateCmd)
	blogsCmd.AddCommand(blogsListCmd)
}
