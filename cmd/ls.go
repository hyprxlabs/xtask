/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"slices"
	"strings"

	"github.com/hyprxlabs/xtask/types"
	"github.com/hyprxlabs/xtask/workflows"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists tasks in the xtaskfile",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		dir, _ := cmd.Flags().GetString("dir")
		file, err := getFile(file, dir)
		if err != nil {
			cmd.PrintErrf("Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		tf := types.NewXTaskfile()
		err = tf.DecodeYAMLFile(file)
		if err != nil {
			cmd.PrintErrf("Error decoding xtaskfile: %v\n", err)
			os.Exit(1)
		}

		wf := workflows.NewWorkflow()
		wf.Args = args
		wf.Context = cmd.Context()

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

		longest := 0
		for _, name := range names {
			if len(name) > longest {
				longest = len(name)
			}
		}
		max := longest + 2

		for _, name := range names {
			desc := ""
			for _, task := range tasks {
				if (task.Name != nil && *task.Name == name) || (task.Name == nil && task.Id == name) {
					if task.Desc != nil && len(*task.Desc) > 0 {
						desc = *task.Desc
					}
					break
				}
			}

			pad := max - len(name)
			if pad < 0 {
				pad = 0
			}

			cmd.Println("\x1b[34m" + name + "\x1b[0m" + strings.Repeat(" ", pad) + "  " + desc)
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
