/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:     "deploy [OPTIONS] [APP...]",
	Short:   "Runs the deploy lifecycle task if it exists",
	Aliases: []string{"up"},
	Long: `Runs the deploy lifecycle task. If no APP is provided, it will run
the default deploy task and the before and after deploy tasks if they exist.

If no apps are provided, the following tasks will be run in order, if they exist.

For example:

Before hook. Runs the first match, if it exists.

- deploy:default:<context>:before
- deploy:<context>:before
- deploy:default:before
- deploy:before

Primary task. Runs the first match, if it exists.

- deploy:default:<context>
- deploy:<context>
- deploy:default
- deploy

After hook. Runs the first match, if it exists.

- deploy:default:<context>:after
- deploy:<context>:after
- deploy:default:after
- deploy:after

If one or more APP are provided, the it will do the following for each app.

With the the context, it will look for tasks in the following 
and execute the first match for the before hook, main task, and 
after hook.

For example, if APP is "web":

Before hook. Runs the first match, if it exists.

- deploy:web:<context>:before
- deploy:web:before
- deploy:before

Primary task. Runs the first match, if it exists.

- deploy:web:<context>
- deploy:web

After hook. Runs the first match, if it exists.

- deploy:web:<context>:after
- deploy:web:after
- deploy:after

If xtask cannot find target with the app name, it will then search
for other xtaskfiles using the app name if the config.dirs.apps section
in the current xtaskfile is set.

If the app name is the same as the last folder in the path. 

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
	Run: func(cmd *cobra.Command, args []string) {
		err := runLifecycle("deploy", cmd)
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	flags := deployCmd.Flags()
	flags.StringArrayP("dotenv", "E", []string{}, "List of dotenv files to load")
	flags.StringToStringP("env", "e", map[string]string{}, "List of environment variables to set")
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
