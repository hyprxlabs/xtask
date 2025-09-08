/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hyprxlabs/xtask/types"
	"github.com/hyprxlabs/xtask/workflows"
	"github.com/spf13/cobra"
)

// manyCmd represents the many command
var manyCmd = &cobra.Command{
	Use:   "many [OPTIONS] [TASK...]",
	Short: "Runs one or more tasks from the xtaskfile",
	Long: `Runs one or more tasks from the xtaskfile. If no TASK is provided, it will run
the default task. No additional arguments are permitted as passing arguments to multiple 
tasks may result in unexpected behavior.
`,
	Run: func(cmd *cobra.Command, args []string) {

		flags := cmd.Flags()

		file, _ := flags.GetString("file")
		dir, _ := flags.GetString("dir")
		targets := flags.Args()
		if len(targets) == 0 {
			targets = []string{"default"}
		}

		file, err := getFile(file, dir)
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

		err = wf.Run(targets, []string{})

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(manyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
