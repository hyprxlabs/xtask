package tasks

import (
	"runtime"
	"strconv"

	"github.com/hyprxlabs/go/exec"

	"github.com/hyprxlabs/xtasks/internal/errors"
	"github.com/hyprxlabs/xtasks/internal/shells"
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

	switch ctx.Task.Uses {
	case "bash":
		cmd = shells.BashScript(ctx.Task.Run)

	case "powershell":
		if runtime.GOOS != "windows" {
			cmd = shells.PwshScript(ctx.Task.Run)
		} else {
			cmd = shells.PowerShellScript(ctx.Task.Run)
		}

	case "sh":
		cmd = shells.ShScriptContext(ctx.Context, ctx.Task.Run)

	case "pwsh":
		cmd = shells.PwshScriptContext(ctx.Context, ctx.Task.Run)

	case "deno":
		cmd = shells.DenoScriptContext(ctx.Context, ctx.Task.Run)

	case "node":
		cmd = shells.NodeScriptContext(ctx.Context, ctx.Task.Run)

	case "bun":
		cmd = shells.BunScriptContext(ctx.Context, ctx.Task.Run)

	case "python":
		cmd = shells.PythonScriptContext(ctx.Context, ctx.Task.Run)

	case "ruby":
		cmd = shells.RubyScriptContext(ctx.Context, ctx.Task.Run)

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
