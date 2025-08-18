package tasks

import (
	"runtime"
	"strconv"

	"github.com/hyprxlabs/go/exec"

	"github.com/hyprxlabs/xtask/internal/errors"
	"github.com/hyprxlabs/xtask/internal/shells"
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
		cmd = shells.BashScriptContext(ctx.Context, run, splat...)

	case "powershell":
		cmd = shells.PowerShellScriptContext(ctx.Context, run, splat...)

	case "sh":
		cmd = shells.ShScriptContext(ctx.Context, run, splat...)

	case "pwsh":
		cmd = shells.PwshScriptContext(ctx.Context, run, splat...)

	case "deno":
		cmd = shells.DenoScriptContext(ctx.Context, run, splat...)

	case "node":
		cmd = shells.NodeScriptContext(ctx.Context, run, splat...)

	case "bun":
		cmd = shells.BunScriptContext(ctx.Context, run, splat...)

	case "python":
		cmd = shells.PythonScriptContext(ctx.Context, run, splat...)

	case "ruby":
		cmd = shells.RubyScriptContext(ctx.Context, run, splat...)

	default:
		err := errors.New("Unsupported shell: " + ctx.Task.Uses)
		return res.Fail(err)
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
