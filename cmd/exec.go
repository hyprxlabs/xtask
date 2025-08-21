/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/internal/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [flags] [command] [args]",
	Short: "executes a command using the environment variables configured in the xtaskfile",
	Long: `Executes a command using the environment variables configured in the xtaskfile
	and does not run any tasks.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, a []string) {
		// task exec

		args := os.Args
		if len(args) > 0 {
			// always will be the cli command
			args = args[1:]
			if len(args) > 0 && args[0] == "exec" {
				args = args[1:]
			} else if len(args) > 0 {
				index := -1
				for i, arg := range args {
					if arg == "exec" {
						index = i
						break
					}
				}

				if index != -1 {
					args = append(args[:index], args[index+1:]...)
				}
			}
		}

		flags := pflag.NewFlagSet("", pflag.ContinueOnError)
		flags.StringP("file", "f", env.Get("XTASK_FILE"), "Path to the xtaskfile (default is ./xtaskfile)")
		flags.StringP("dir", "d", env.Get("XTASK_DIR"), "Directory to run the task in (default is current directory)")
		flags.StringArrayP("dotenv", "E", []string{}, "List of dotenv files to load")
		flags.StringToStringP("env", "e", map[string]string{}, "List of environment variables to set")

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

		file, _ := flags.GetString("file")
		dir, _ := flags.GetString("dir")
		file, err = getFile(file, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		if len(remainingArgs) == 0 {
			os.Stderr.WriteString("No command provided to exec.\n")
			cmd.Help()
			os.Exit(1)
		}

		dotenvFiles, _ := flags.GetStringArray("dotenv")
		envMap, _ := flags.GetStringToString("env")

		err = workflow.Run(workflow.Params{
			Args:                remainingArgs,
			Tasks:               []string{"default"},
			Timeout:             0,
			CommandSubstitution: true,
			Context:             cmd.Context(),
			Command:             "exec",
			File:                file,
			Dotenv:              dotenvFiles,
			Env:                 envMap,
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
	flags.StringArrayP("dotenv", "E", []string{}, "List of dotenv files to load")
	flags.StringToStringP("env", "e", nil, "Environment variables to set for the command")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
