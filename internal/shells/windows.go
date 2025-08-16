//go:build windows

package shells

import (
	"context"
	"path/filepath"
	"strings"

	e "os/exec"

	"github.com/hyprxlabs/go/exec"
)

func init() {
	exec.Register("bash", &exec.Executable{
		Name:     "bash",
		Variable: "BASH_PATH",
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\bash.exe",
			"%ProgramFiles(x86)%\\Git\\bin\\bash.exe",
			"${SystemRoot}\\System32\\bash.exe",
		},
	})

	exec.Register("pwsh", &exec.Executable{
		Name:     "pwsh",
		Variable: "PWSH_PATH",
		Windows: []string{
			"${ProgramFiles}\\PowerShell\\7\\pwsh.exe",
			"${ProgramFiles(x86)}\\PowerShell\\7\\pwsh.exe",
			"${ProgramFileds}\\PowerShell\\6\\pwsh.exe",
			"${ProgramFiles(x86)}\\PowerShell\\6\\pwsh.exe",
		},
	})

	exec.Register("powershell", &exec.Executable{
		Name:     "powershell",
		Variable: "POWERSHELL_PATH",
		Windows: []string{
			"${SystemRoot}\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
			"${SystemRoot}\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe",
		},
	})

	exec.Register("sh", &exec.Executable{
		Name:     "sh",
		Variable: "SH_PATH",
		Windows: []string{
			"${ProgramFiles}\\Git\\bin\\sh.exe",
			"%ProgramFiles(x86)%\\Git\\bin\\sh.exe",
		},
	})

	exec.Register("deno", &exec.Executable{
		Name:     "deno",
		Variable: "DENO_PATH",
		Windows: []string{
			"${USERPROFILE}\\.deno\\bin\\deno.exe",
			"${LOCALAPPDATA}\\Programs\\bin\\deno.exe",
			"${LOCALAPPDATA}\\Microsoft\\WinGet\\Packages\\DenoLand.Deno_Microsoft.Winget.Source_8wekyb3d8bbwe\\deno.exe",
			"${ProgramFiles}\\Deno\\deno.exe",
		},
	})

	exec.Register("node", &exec.Executable{
		Name:     "node",
		Variable: "NODE_PATH",
		Windows: []string{
			"${ProgramFiles}\\nodejs\\node.exe",
			"${ProgramFiles(x86)}\\nodejs\\node.exe",
			"${ProgramFiles}\\node\\node.exe",
			"${ProgramFiles(x86)}\\node\\node.exe",
		},
	})

	exec.Register("bun", &exec.Executable{
		Name:     "bun",
		Variable: "BUN_PATH",
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
		Variable: "PYTHON_PATH",
		Windows: []string{
			"${ProgramFiles}\\Python\\Python.exe",
			"${ProgramFiles(x86)}\\Python\\Python.exe",
		},
	})

	exec.Register("ruby", &exec.Executable{
		Name:     "ruby",
		Variable: "RUBY_PATH",
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
	bash, _ := exec.Which("bash")
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

func BashScriptContext(ctx context.Context, script string) *exec.Cmd {

	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".sh") {
		script = resolveScriptFile(script)
	}

	args := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}

	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "bash", args...),
	}

	return cmd
}
