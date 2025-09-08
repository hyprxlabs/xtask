/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/types"
	"github.com/hyprxlabs/xtask/workflows"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [OPTIONS] [TASK...] [--] [REMAINING_ARGS...]",
	Short: "Runs a single task from the xtaskfile and may pass remaining arguments to it.",
	Long: `Run a single task from the xtaskfile.
Additional arguments may be passed to the task. The -- separator is may be used to 
force all subsequent arguments to be treated as remaining arguments.`,
	Example: `xtask run test
  xtask run -c CONTEXTA -e MY_VAR=test build -- --no-cache
  xtask run -e ENV=production deploy --tag v1.0.0`,
	Aliases:            []string{"r"},
	Args:               cobra.ArbitraryArgs,
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
		flags.StringP("context", "c", env.Get("XTASK_CONTEXT"), "Context to use.")

		targets := []string{}
		cmdArgs := []string{}
		remainingArgs := []string{}
		size := len(args)
		inRemaining := false
		for i := 0; i < size; i++ {
			n := args[i]
			if n == "--" {
				inRemaining = true
				continue
			}

			if inRemaining {
				remainingArgs = append(remainingArgs, args[i])
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

			targets = append(targets, n)
			inRemaining = true
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

		tf := types.NewXTaskfile()

		err = tf.DecodeYAMLFile(file)
		tf.Path = file

		if err != nil {
			cmd.PrintErrf("Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		if len(dotenvFiles) > 0 {
			tf.Dotenv = append(tf.Dotenv, dotenvFiles...)
		}

		if len(envVars) > 0 {
			if tf.Env == nil {
				tf.Env = types.NewEnv()
			}

			for k, v := range envVars {
				tf.Env.Set(k, v)
			}
		}

		wf := workflows.NewWorkflow()

		err = wf.Load(*tf)
		if err != nil {
			cmd.PrintErrf("Error loading xtaskfile: %v\n", err)
			os.Exit(1)
		}

		err = wf.Run(targets, remainingArgs)

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
	runCmd.Flags().StringToStringP("env", "e", map[string]string{}, "List of environment variables to  ")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
