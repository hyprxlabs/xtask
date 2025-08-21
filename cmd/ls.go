/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hyprxlabs/xtask/internal/workflow"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "lists tasks in the xtaskfile",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		dir, _ := cmd.Flags().GetString("dir")
		file, err := getFile(file, dir)
		if err != nil {
			cmd.PrintErrf("Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		err = workflow.Run(workflow.Params{
			Args:                args,
			Tasks:               []string{},
			Timeout:             0,
			CommandSubstitution: true,
			Context:             cmd.Context(),
			File:                file,
			Command:             "ls",
		})

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
