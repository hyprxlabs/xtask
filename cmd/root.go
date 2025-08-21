/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hyprxlabs/go/env"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "task",
	Short:   "a cross platform task runner",
	Version: Version,
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	args := os.Args
	commands := []string{"run", "ls", "exec"}
	if len(args) == 1 {
		rootCmd.SetArgs([]string{"run"})
	} else {
		if len(args) == 2 {
			flag := args[1]
			if flag == "--help" || flag == "-h" || flag == "-v" || flag == "--version" {
				rootCmd.SetArgs([]string{flag})
			} else {
				validCommand := false
				for _, cmd := range commands {
					if flag == cmd {
						rootCmd.SetArgs([]string{cmd})
						validCommand = true
						break
					}
				}

				if !validCommand {
					rootCmd.SetArgs([]string{"run", flag})
				}
			}
		} else {

			hasCommand := false

			for i := 0; i < len(args); i++ {
				arg := args[i]
				switch arg {
				case "--file", "-f":
					i++

				case "--dir", "-d":
					i++
				case "xtask":
					continue
				case "ls":
					hasCommand = true
				case "run":
					hasCommand = true
				case "exec":
					hasCommand = true
				}
			}
			if !hasCommand {
				nargs := []string{"run"}
				nargs = append(nargs, args[1:]...)
				rootCmd.SetArgs(nargs)
			}
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	file := env.Get("XTASK_FILE")
	dir := env.Get("XTASK_DIR")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("file", "f", file, "Path to the YAML file")
	rootCmd.PersistentFlags().StringP("dir", "d", dir, "Directory to run the task in (default is current directory)")
}
