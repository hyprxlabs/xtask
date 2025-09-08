/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"slices"
	"strings"

	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/versions"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "xtask",
	Short:   "a cross platform task runner",
	Version: versions.Version,
	Long: `A cross platform task runner
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	args := os.Args
	commands := []string{
		"audit",
		"b",
		"build",
		"completion",
		"deploy",
		"destroy",
		"down",
		"exec",
		"help",
		"install",
		"many",
		"pack",
		"publish",
		"ls",
		"lc",
		"lifecycle",
		"run",
		"runlc",
		"test",
		"t",
		"uninstall",
		"upgrade",
		"up",
		"version",
		"x"}

	hasCommand := false
	hasPossibleTarget := false

	for _, arg := range args {
		if slices.Contains(commands, arg) {
			hasCommand = true
		}

		if !strings.HasPrefix(arg, "-") {
			hasPossibleTarget = true
		}
	}

	if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" || args[1] == "help" {
			rootCmd.Help()
			os.Exit(0)
		}

		if args[1] == "-v" || args[1] == "--version" || args[1] == "version" {
			os.Stdout.WriteString("xtask version " + versions.Version + "\n")
			os.Exit(0)
		}
	}

	if !hasCommand && hasPossibleTarget {
		args = append([]string{args[0], "run"}, args[1:]...)
		os.Args = args

		os.Stdout.WriteString(strings.Join(os.Args[1:], " ") + "\n")
		rootCmd.SetArgs(args[1:])
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	file := env.Get("XTASK_FILE")
	dir := env.Get("XTASK_DIR")
	context := env.Get("XTASK_CONTEXT")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("file", "f", file, "Path to the YAML file.")
	rootCmd.PersistentFlags().StringP("dir", "d", dir, "Directory to run the task in (default is current directory).")
	rootCmd.PersistentFlags().StringP("context", "c", context, "The context to use. If not set, the 'default' context is used.")
}
