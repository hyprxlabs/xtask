/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"slices"

	"github.com/hyprxlabs/xtask/types"
	"github.com/hyprxlabs/xtask/workflows"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

		wf := workflows.NewWorkflow()
		wf.Args = args
		wf.Context = cmd.Context()

		data, err := os.ReadFile(file)
		if err != nil {
			cmd.PrintErrf("Error reading xtaskfile: %v\n", err)
			os.Exit(1)
		}
		tf := &types.XTaskfile{}
		err = yaml.Unmarshal(data, tf)
		if err != nil {
			cmd.PrintErrf("Error parsing xtaskfile: %v\n", err)
			os.Exit(1)
		}

		err = wf.Load(*tf)
		if err != nil {
			cmd.PrintErrf("Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		tasks := wf.List()

		names := []string{}
		for _, task := range tasks {
			if task.Name != nil && len(*task.Name) > 0 {
				names = append(names, *task.Name)
			} else {
				names = append(names, task.Id)
			}
		}
		slices.Sort(names)

		for _, name := range names {
			cmd.Println(name)
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
