//go:build !windows

package shells

import (
	"context"
	"strings"

	"github.com/hyprxlabs/go/cmdargs"
	"github.com/hyprxlabs/go/exec"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: "XTASK_BASH_EXE",
		Linux: []string{
			"/bin/bash",
			"/usr/bin/bash",
		},
	})

	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: "XTASK_PWSH_EXE",
		Linux: []string{
			"/usr/bin/pwsh",
			"/usr/local/bin/pwsh",
		},
	})

	exec.Register("powershell", &exec.Executable{
		Name:     "powershell",
		Variable: "XTASK_POWERSHELL_EXE",
		Linux: []string{
			"/usr/bin/pwsh",
			"/usr/local/bin/pwsh",
		},
	})

	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: "XTASK_SH_EXE",
		Linux: []string{
			"/bin/sh",
			"/usr/bin/sh",
		},
	})

	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: "XTASK_DENO_EXE",
		Linux: []string{
			"${HOME}/.local/bin/deno",
			"${HOME}/.deno/bin/deno",
			"/usr/bin/deno",
			"/usr/local/bin/deno",
		},
	})

	exec.Register("node", &exec.Executable{
		Name:     "node",
		Variable: "XTASK_NODE_EXE",
		Linux: []string{
			"/usr/bin/node",
			"/usr/local/bin/node",
		},
	})

	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: "XTASK_BUN_EXE",
		Linux: []string{
			"/usr/bin/bun",
			"/usr/local/bin/bun",
		},
	})

	exec.Register("python", &exec.Executable{
		Name:     "python",
		Variable: "XTASK_PYTHON_EXE",
		Linux: []string{
			"/usr/bin/python",
			"/usr/bin/python3",
			"/usr/local/bin/python",
			"/usr/local/bin/python3",
		},
	})

	exec.Register("ruby", &exec.Executable{
		Name:     "ruby",
		Variable: "XTASK_RUBY_EXE",
		Linux: []string{
			"/usr/bin/ruby",
			"/usr/local/bin/ruby",
		},
	})

}

func BashScript(script string) *exec.Cmd {
	args := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	return exec.New("bash", args...)
}

func BashScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	noLines := !strings.ContainsAny(script, "\n\r")

	if len(args) > 0 && noLines {
		next := cmdargs.New([]string{script}).Append(args...).String()
		script = next
	}

	splat := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	exe, _ := exec.Find("bash", nil)
	if exe == "" {
		exe = "bash"
	}
	return exec.NewContext(ctx, exe, splat...)
}
