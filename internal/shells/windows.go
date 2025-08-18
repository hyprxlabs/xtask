//go:build windows

package shells

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hyprxlabs/go/cmdargs"
	"github.com/hyprxlabs/go/exec"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: "XTASK_WIN_BASH_EXE",
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\bash.exe",
			"%ProgramFiles(x86)%\\Git\\bin\\bash.exe",
			"${SystemRoot}\\System32\\bash.exe",
		},
	})

	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: "XTASK_WIN_PWSH_EXE",
		Windows: []string{
			"${ProgramFiles}\\PowerShell\\7\\pwsh.exe",
			"${ProgramFiles(x86)}\\PowerShell\\7\\pwsh.exe",
			"${ProgramFileds}\\PowerShell\\6\\pwsh.exe",
			"${ProgramFiles(x86)}\\PowerShell\\6\\pwsh.exe",
		},
	})

	exec.Register("powershell", &exec.Executable{
		Name:     "powershell",
		Variable: "XTASK_WIN_POWERSHELL_EXE",
		Windows: []string{
			"${SystemRoot}\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
			"${SystemRoot}\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe",
		},
	})

	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: "XTASK_WIN_SH_EXE",
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\sh.exe",
			"%ProgramFiles(x86)%\\Git\\bin\\sh.exe",
		},
	})

	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: "XTASK_WIN_DENO_EXE",
		Windows: []string{
			"${USERPROFILE}\\.deno\\bin\\deno.exe",
			"${LOCALAPPDATA}\\Programs\\bin\\deno.exe",
			"${LOCALAPPDATA}\\Microsoft\\WinGet\\Packages\\DenoLand.Deno_Microsoft.Winget.Source_8wekyb3d8bbwe\\deno.exe",
			"${ProgramFiles}\\Deno\\deno.exe",
		},
	})

	exec.Register("node", &exec.Executable{
		Name:     "node",
		Variable: "XTASK_WIN_NODE_EXE",
		Windows: []string{
			"${ProgramFiles}\\nodejs\\node.exe",
			"${ProgramFiles(x86)}\\nodejs\\node.exe",
			"${ProgramFiles}\\node\\node.exe",
			"${ProgramFiles(x86)}\\node\\node.exe",
		},
	})

	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: "XTASK_WIN_BUN_EXE",
		Windows: []string{
			"${USERPROFILE}\\.bun\\bin\\bun.exe",
			"${LOCALAPPDATA}\\Programs\\bin\\bun.exe",
			"${LOCALAPPDATA}\\Microsoft\\WinGet\\Links\\bin.exe",
			"${ProgramFiles}\\bun\\bin\\bun.exe",
			"${ProgramFiles(x86)}\\bun\\bin\\bun.exe",
		},
	})

	exec.Register("python", &exec.Executable{
		Name:     "python",
		Variable: "XTASK_WIN_PYTHON_EXE",
		Windows: []string{
			"${ProgramFiles}\\Python\\Python.exe",
			"${ProgramFiles(x86)}\\Python\\Python.exe",
		},
	})

	exec.Register("ruby", &exec.Executable{
		Name:     "ruby",
		Variable: "XTASK_WIN_RUBY_EXE",
		Windows: []string{
			"${ProgramFiles}\\Ruby\\bin\\ruby.exe",
			"${ProgramFiles(x86)}\\Ruby\\bin\\ruby.exe",
		},
	})
}

func resolveScriptFile(script string) string {
	if !filepath.IsAbs(script) {
		file, err := filepath.Abs(script)
		if err != nil {
			script = file
		}
	}

	// determine if bash is the WSL one.
	bash, _ := exec.Find("bash", nil)
	if !strings.Contains(strings.ToLower(bash), "system32") {
		return script
	}

	script = strings.ReplaceAll(script, "\\", "/")
	if script[1] == ':' {
		script = "/mnt/" + strings.ToLower(script[0:1]) + script[2:]
	}

	return script
}

func BashScript(script string) *exec.Cmd {
	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".sh") {
		script = resolveScriptFile(script)
	}

	args := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	return exec.New("bash", args...)
}

func BashScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	noLines := !strings.ContainsAny(script, "\n\r")
	exe, _ := exec.Find("bash", nil)
	if exe == "" {
		exe = "bash"
	}

	if noLines && strings.HasSuffix(strings.TrimSpace(script), ".sh") {
		script = resolveScriptFile(script)
	}

	if len(args) > 0 && noLines {
		next := cmdargs.New([]string{script}).Append(args...).String()
		script = next
	}

	splat := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	return exec.NewContext(ctx, exe, splat...)
}
