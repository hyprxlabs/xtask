/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// runlcCmd represents the runlc command
var runlcCmd = &cobra.Command{
	Use:     "runlc [OPTIONS] <task> [APP...]",
	Short:   "Runs a lifecycle task if it exists",
	Aliases: []string{"lc", "lifecycle"},
	Long: `Runs the given lifecycle <task>.

The task name is required.

If no apps are provided, the following tasks will be run in order, if they exist.

For example:

Before hook. Runs the first match, if it exists.

- <task>:default:<context>:before
- <task>:<context>:before
- <task>:default:before
- <task>:before

Primary task. Runs the first match, if it exists.

- <task>:default:<context>
- <task>:<context>
- <task>:default
- <task>

After hook. Runs the first match, if it exists.

- <task>:default:<context>:after
- <task>:<context>:after
- <task>:default:after
- <task>:after

If one or more APP are provided, the it will do the following for each app.

With the context, it will look for tasks in the following 
and execute the first match for the before hook, main task, and 
after hook.

For example, if APP is "web":

Before hook. Runs the first match, if it exists.

- <task>:web:<context>:before
- <task>:web:before
- <task>:before

Primary task. Runs the first match, if it exists.

- <task>:web:<context>
- <task>:web

After hook. Runs the first match, if it exists.

- <task>:web:<context>:after
- <task>:web:after
- <task>:after

If xtask cannot find target with the app name, it will then search
other file system directories for xtaskfiles if there are directories
configured in  "config.dirs.apps" section of the current xtaskfile.

If the app name is the same as the last folder in the path. If the for
example the app name is "web" and the current path is:

./src/web

if it will first for the following in order...

./src/web/[CONTEXT]/xtaskfile
./src/web/xtaskfile

If the app name is not the same as the last folder in the path e.g.

./src

Then it will search for...

./src/[APP]/[CONTEXT]/xtaskfile
./src/[APP]/xtaskfile

If no context is provided, it will use "default" as the context.

You can override the context by setting the XTASK_CONTEXT environment
variable or using the --context, -c flag`,
	Example: `xtask runlc -c Release build app1 app2
xtask lifecycle -c production deploy
xtask lc -c production deploy app1 app2
	`,
	Run: func(cmd *cobra.Command, a []string) {

		args := os.Args[1:]
		target := ""

		index := -1
		for i, arg := range args {
			if arg == "runlc" {
				index = i
				break
			}
		}

		if index != -1 {
			if index+1 >= len(args) {
				fmt.Println("Error: target is required")
				os.Exit(1)
			}

			target = args[index+1]
		}

		err := runLifecycle(target, cmd)
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	flags := runlcCmd.Flags()
	flags.StringArrayP("dotenv", "E", []string{}, "List of dotenv files to load")
	flags.StringToStringP("env", "e", map[string]string{}, "List of environment variables to set")
	rootCmd.AddCommand(runlcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runlcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runlcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
