/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/internal/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, a []string) {

		args := os.Args

		if len(args) > 0 {
			// always will be the cli command
			args = args[1:]

			if len(args) > 0 && args[0] == "run" {
				args = args[1:]
			} else if len(args) > 0 {
				index := -1
				for i, arg := range args {
					if arg == "run" {
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

		targets := []string{}
		cmdArgs := []string{}
		remainingArgs := []string{}
		size := len(args)
		inRemaining := false
		inTargets := false
		for i := 0; i < size; i++ {
			n := args[i]
			if n == "--" {
				inTargets = false
				inRemaining = true
				continue
			}

			if inRemaining {
				inTargets = false
				remainingArgs = append(remainingArgs, args[i])
				continue
			}

			if inTargets {
				if n == "--" {
					inTargets = false
					inRemaining = true
					continue
				}

				if len(n) > 0 && n[0] == '-' {
					inTargets = false
					inRemaining = true
					remainingArgs = append(remainingArgs, n)
					continue
				}

				targets = append(targets, args[i])
				continue
			}

			if len(n) > 0 && n[0] == '-' {
				cmdArgs = append(cmdArgs, n)
				j := i + 1
				if j < size && len(args[j]) > 0 && args[j][0] != '-' {
					cmdArgs = append(cmdArgs, args[j])
					i++ // Skip the next argument as it's a value for the flag
				}

				continue
			}

			inTargets = true
			targets = append(targets, n)
		}

		if len(targets) == 0 {
			targets = append(targets, "default")
		}

		err := flags.Parse(cmdArgs)
		if err != nil {
			cmd.PrintErrf("Error parsing flags: %v\n", err)
			os.Exit(1)
		}

		file, _ := flags.GetString("file")
		dir, _ := flags.GetString("dir")

		file, err = getFile(file, dir)
		if err != nil {
			cmd.PrintErrf("Error resolving file: %v\n", err)
			os.Exit(1)
		}

		dotenvFiles, _ := flags.GetStringArray("dotenv")
		envVars, _ := flags.GetStringToString("env")

		err = workflow.Run(workflow.Params{
			Args:                remainingArgs,
			Tasks:               targets,
			Timeout:             0,
			CommandSubstitution: true,
			Context:             cmd.Context(),
			Command:             "run",
			File:                file,
			Dotenv:              dotenvFiles,
			Env:                 envVars,
		})

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringArrayP("dotenv", "E", []string{}, "List of dotenv files to load")
	runCmd.Flags().StringToStringP("env", "e", map[string]string{}, "List of environment variables to set")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
