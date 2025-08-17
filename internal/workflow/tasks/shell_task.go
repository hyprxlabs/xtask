package tasks

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/hyprxlabs/go/cmdargs"
	"github.com/hyprxlabs/go/exec"

	"github.com/hyprxlabs/task/internal/errors"
	"github.com/hyprxlabs/task/internal/shells"
)

func runShell(ctx TaskContext) *TaskResult {
	res := NewTaskResult()
	if ctx.Task.Uses == "" {
		shell := "bash"
		if runtime.GOOS == "windows" {
			shell = "powershell"
		}

		ctx.Task.Uses = shell
	}

	var cmd *exec.Cmd

	run := ctx.Task.Run
	splat := ctx.Task.Args

	switch ctx.Task.Uses {
	case "bash":
		if len(splat) > 0 {
			r := strings.TrimSpace(run)
			nargs := cmdargs.New([]string{r})
			nargs.Append(splat...)
			run = nargs.String()
		}

		cmd = shells.BashScriptContext(ctx.Context, run)

	case "powershell":
		if len(splat) > 0 {
			r := strings.TrimSpace(run)
			nargs := cmdargs.New([]string{r})
			nargs.Append(splat...)
			run = nargs.String()
		}
		if runtime.GOOS != "windows" {
			cmd = shells.PwshScriptContext(ctx.Context, run)
		} else {
			cmd = shells.PowerShellScriptContext(ctx.Context, run)
		}

	case "sh":
		if len(splat) > 0 {
			r := strings.TrimSpace(run)
			nargs := cmdargs.New([]string{r})
			nargs.Append(splat...)
			run = nargs.String()
		}
		cmd = shells.ShScriptContext(ctx.Context, run)

	case "pwsh":
		if len(splat) > 0 {
			r := strings.TrimSpace(run)
			nargs := cmdargs.New([]string{r})
			nargs.Append(splat...)
			run = nargs.String()
		}
		cmd = shells.PwshScriptContext(ctx.Context, run)

	case "deno":
		cmd = shells.DenoScriptContext(ctx.Context, run)

	case "node":
		cmd = shells.NodeScriptContext(ctx.Context, run)

	case "bun":
		cmd = shells.BunScriptContext(ctx.Context, run)

	case "python":
		cmd = shells.PythonScriptContext(ctx.Context, run)

	case "ruby":
		cmd = shells.RubyScriptContext(ctx.Context, run)

	default:
		err := errors.New("Unsupported shell: " + ctx.Task.Uses)
		return res.Fail(err)
	}

	if ctx.Task.Timeout > 0 {
		cmd.WithTimeout(ctx.Task.Timeout)
	}

	if ctx.Task.Cwd != "" {
		cmd.Dir = ctx.Task.Cwd
	}

	if len(ctx.Task.Env) > 0 {
		cmd.WithEnvMap(ctx.Task.Env)
	}

	res.Start()
	o, err := cmd.Run()
	if err != nil {
		return res.Fail(err)
	}

	if o.Code != 0 {
		err := errors.New("Task " + ctx.Task.Id + " failed with exit code " + strconv.Itoa(o.Code))
		return res.Fail(err)
	}

	// Placeholder for running a shell command
	// This would typically involve executing the command in the shell
	return res.Ok()
}
