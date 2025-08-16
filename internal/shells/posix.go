//go:build !windows

package shells

import (
	"context"

	e "os/exec"

	"github.com/hyprxlabs/go/exec"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: "BASH_PATH",
		Linux: []string{
			"/bin/bash",
			"/usr/bin/bash",
		},
	})

	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: "PWSH_PATH",
		Linux: []string{
			"/usr/bin/pwsh",
			"/usr/local/bin/pwsh",
		},
	})

	exec.Register("powershell", &exec.Executable{
		Name:     "powershell",
		Variable: "POWERSHELL_PATH",
		Linux: []string{
			"/usr/bin/pwsh",
			"/usr/local/bin/pwsh",
		},
	})

	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: "SH_PATH",
		Linux: []string{
			"/bin/sh",
			"/usr/bin/sh",
		},
	})

	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: "DENO_PATH",
		Linux: []string{
			"${HOME}/.local/bin/deno",
			"${HOME}/.deno/bin/deno",
			"/usr/bin/deno",
			"/usr/local/bin/deno",
		},
	})

	exec.Register("node", &exec.Executable{
		Name:     "node",
		Variable: "NODE_PATH",
		Linux: []string{
			"/usr/bin/node",
			"/usr/local/bin/node",
		},
	})

	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: "BUN_PATH",
		Linux: []string{
			"/usr/bin/bun",
			"/usr/local/bin/bun",
		},
	})

	exec.Register("python", &exec.Executable{
		Name:     "python",
		Variable: "PYTHON_PATH",
		Linux: []string{
			"/usr/bin/python",
			"/usr/bin/python3",
			"/usr/local/bin/python",
			"/usr/local/bin/python3",
		},
	})

	exec.Register("ruby", &exec.Executable{
		Name:     "ruby",
		Variable: "RUBY_PATH",
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

func BashScriptContext(ctx context.Context, script string) *exec.Cmd {
	args := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "bash", args...),
	}
	return cmd

}
