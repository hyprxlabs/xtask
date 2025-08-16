package shells

import (
	"context"
	e "os/exec"
	"strings"

	"github.com/hyprxlabs/go/exec"
)

func PwshScript(script string) *exec.Cmd {
	args := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}

	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		args = append(args, "-File", strings.TrimSpace(script))
	} else {
		args = append(args, "-Command", script)
	}

	return exec.New("pwsh", args...)
}

func PwshScriptContext(ctx context.Context, script string) *exec.Cmd {
	args := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}

	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		args = append(args, "-File", strings.TrimSpace(script))
	} else {
		args = append(args, "-Command", script)
	}

	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "pwsh", args...),
	}

	return cmd
}

func PowerShellScript(script string) *exec.Cmd {
	args := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}

	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		args = append(args, "-File", strings.TrimSpace(script))
	} else {
		args = append(args, "-Command", script)
	}

	return exec.New("powershell", args...)
}

func PowerShellScriptContext(ctx context.Context, script string) *exec.Cmd {
	args := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}

	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if (!strings.ContainsAny(script, "\n\r")) && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		args = append(args, "-File", strings.TrimSpace(script))
	} else {
		args = append(args, "-Command", script)
	}

	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "powershell", args...),
	}

	return cmd
}

func ShScript(script string) *exec.Cmd {
	args := []string{"-e", script}

	return exec.New("sh", args...)
}

func ShScriptContext(ctx context.Context, script string) *exec.Cmd {
	args := []string{"-e", script}

	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "sh", args...),
	}

	return cmd
}

func DenoScript(script string) *exec.Cmd {

	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".ts") || strings.HasSuffix(trimmed, ".js") {
			return exec.New("deno", "run", "-A", "--allow-scripts", trimmed)
		}
	}

	args := []string{"eval", "--ext=ts", "--allow-scripts", script}
	return exec.New("deno", args...)
}

func DenoScriptContext(ctx context.Context, script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".ts") || strings.HasSuffix(trimmed, ".js") {
			return exec.New("deno", "run", "-A", "--allow-scripts", trimmed)
		}
	}

	args := []string{"eval", "--ext=ts", "--allow-scripts", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "deno", args...),
	}

	return cmd
}

func NodeScript(script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") {
			return exec.New("node", trimmed)
		}

		if strings.HasSuffix(trimmed, ".ts") {
			return exec.New("node", "--experimental-transform-types", trimmed)
		}
	}

	args := []string{"--experimental-transform-types", "-e", script}
	return exec.New("node", args...)
}

func NodeScriptContext(ctx context.Context, script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") {
			return exec.New("node", trimmed)
		}

		if strings.HasSuffix(trimmed, ".ts") {
			return exec.New("node", "--experimental-transform-types", trimmed)
		}
	}

	args := []string{"--experimental-transform-types", "-e", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "node", args...),
	}

	return cmd
}

func BunScript(script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") || strings.HasSuffix(trimmed, ".ts") {
			return exec.New("bun", "run", trimmed)
		}
	}

	args := []string{"-e", script}
	return exec.New("bun", args...)
}

func BunScriptContext(ctx context.Context, script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") || strings.HasSuffix(trimmed, ".ts") {
			return exec.New("bun", "run", trimmed)
		}
	}

	args := []string{"-e", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "bun", args...),
	}

	return cmd
}

func PythonScript(script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".py") {
			return exec.New("python", trimmed)
		}
	}

	args := []string{"-c", script}
	return exec.New("python", args...)
}

func PythonScriptContext(ctx context.Context, script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".py") {
			return exec.New("python", trimmed)
		}
	}

	args := []string{"-c", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "python", args...),
	}

	return cmd
}

func RubyScript(script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".rb") {
			return exec.New("ruby", trimmed)
		}
	}

	args := []string{"-e", script}
	return exec.New("ruby", args...)
}

func RubyScriptContext(ctx context.Context, script string) *exec.Cmd {
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".rb") {
			return exec.New("ruby", trimmed)
		}
	}

	args := []string{"-e", script}
	cmd := &exec.Cmd{
		Cmd: e.CommandContext(ctx, "ruby", args...),
	}

	return cmd
}

type ScriptAttributes struct {
	Env    map[string]string
	Dir    string
	Output bool
}

type commandArgs struct {
	Args []string
}

type tokenGroup struct {
	Commands []commandArgs
	Kind     string
}
