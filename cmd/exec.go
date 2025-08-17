/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hyprxlabs/task/internal/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "executes a command using the environment variables configured in the xtaskfile",
	Long: `Executes a command using the environment variables configured in the xtaskfile
	and does not run any tasks.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, a []string) {
		// task exec
		args := os.Args[2:]
		flags := pflag.NewFlagSet("", pflag.ContinueOnError)

		cmdArgs := []string{}
		remainingArgs := []string{}
		size := len(args)
		inRemaining := false
		for i := 0; i < size; i++ {
			if inRemaining {
				remainingArgs = append(remainingArgs, args[i])
				continue
			}

			n := args[i]
			if len(n) > 0 && n[0] == '-' {
				cmdArgs = append(cmdArgs, n)
				j := i + 1
				if j < size && len(args[j]) > 0 && args[j][0] != '-' {
					cmdArgs = append(cmdArgs, args[j])
					i++ // Skip the next argument as it's a value for the flag
				}

				continue
			}

			inRemaining = true
			remainingArgs = append(remainingArgs, n)
		}

		err := flags.Parse(cmdArgs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}

		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			file = "./xtaskfile"
		}

		err = workflow.Run(workflow.Params{
			Args:                remainingArgs,
			Tasks:               []string{"default"},
			Timeout:             0,
			CommandSubstitution: true,
			Context:             cmd.Context(),
			Command:             "exec",
			File:                file,
		})

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	flags := execCmd.Flags()
	flags.StringP("file", "f", "", "Path to the xtaskfile (default is ./xtaskfile)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
