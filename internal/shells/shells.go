package shells

import (
	"context"
	"strings"

	"github.com/hyprxlabs/go/cmdargs"
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

func PwshScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	splat := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}
	noLines := !strings.ContainsAny(script, "\n\r")
	exe, _ := exec.Find("pwsh", nil)
	if exe == "" {
		exe = "pwsh"
	}

	// if script is a single line and ends with .ps1, use -File, otherwise use -Command
	// for single line scripts, allow appending additional arguments
	if noLines && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		if len(args) > 0 {
			next := cmdargs.New([]string{script}).Append(args...).String()
			script = next
		}
		splat = append(splat, "-File", strings.TrimSpace(script))
	} else {
		if noLines && len(args) > 0 {
			next := cmdargs.New([]string{script}).Append(args...).String()
			script = next
		}
		splat = append(splat, "-Command", script)
	}

	return exec.NewContext(ctx, exe, splat...)
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

func PowerShellScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("powershell", nil)
	if exe == "" {
		exe = "powershell"
	}
	splat := []string{"-NoProfile", "-NoLogo", "-ExecutionPolicy", "Bypass"}

	noLines := !strings.ContainsAny(script, "\n\r")
	// if script is a single line and ends with .sh, resolve it to an absolute path
	// and for windows, convert it to a WSL path if necessary
	if noLines && strings.HasSuffix(strings.TrimSpace(script), ".ps1") {
		if len(args) > 0 {
			next := cmdargs.New([]string{script}).Append(args...).String()
			script = next
		}

		splat = append(splat, "-File", strings.TrimSpace(script))
	} else {
		if noLines && len(args) > 0 {
			next := cmdargs.New([]string{script}).Append(args...).String()
			script = next
		}

		splat = append(splat, "-Command", script)
	}

	return exec.NewContext(ctx, exe, splat...)
}

func ShScript(script string) *exec.Cmd {
	args := []string{"-e", script}

	return exec.New("sh", args...)
}

func ShScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("sh", nil)
	if exe == "" {
		exe = "sh"
	}
	if len(args) > 0 && !strings.ContainsAny(script, "\n\r") {
		next := cmdargs.New([]string{script}).Append(args...).String()
		script = next
	}

	splat := []string{"-e", script}
	return exec.NewContext(ctx, exe, splat...)
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

func DenoScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("deno", nil)
	if exe == "" {
		exe = "deno"
	}

	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".ts") || strings.HasSuffix(trimmed, ".js") {
			splat := []string{"run", "-A", "--allow-scripts", trimmed}
			if len(args) > 0 {
				splat = append(splat, args...)
			}
			return exec.NewContext(ctx, exe, splat...)
		}
	}

	splat := []string{"eval", "--ext=ts", "--allow-scripts", script}
	return exec.NewContext(ctx, exe, splat...)
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

func NodeScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("node", nil)
	if exe == "" {
		exe = "node"
	}

	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") {
			if len(args) > 0 {
				next := cmdargs.New([]string{trimmed}).Append(args...).String()
				return exec.NewContext(ctx, exe, next)
			}

			return exec.NewContext(ctx, exe, trimmed)
		}

		if strings.HasSuffix(trimmed, ".ts") {
			if len(args) > 0 {
				next := cmdargs.New([]string{trimmed}).Append(args...).String()
				return exec.NewContext(ctx, exe, "--experimental-transform-types", next)
			}

			return exec.NewContext(ctx, exe, "--experimental-transform-types", trimmed)
		}
	}

	splat := []string{"-e", script}
	return exec.NewContext(ctx, exe, splat...)
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

func BunScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("bun", nil)
	if exe == "" {
		exe = "bun"
	}
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".js") || strings.HasSuffix(trimmed, ".mjs") || strings.HasSuffix(trimmed, ".ts") {
			if len(args) > 0 {
				next := cmdargs.New([]string{trimmed}).Append(args...).String()
				return exec.NewContext(ctx, exe, "run", next)
			}
			return exec.NewContext(ctx, exe, "run", trimmed)
		}
	}

	splat := []string{"-e", script}
	return exec.NewContext(ctx, exe, splat...)
}

func PythonScript(script string) *exec.Cmd {
	exe, _ := exec.Find("python", nil)
	if exe == "" {
		exe = "python"
	}

	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".py") {
			return exec.New(exe, trimmed)
		}
	}

	args := []string{"-c", script}
	return exec.New(exe, args...)
}

func PythonScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("python", nil)
	if exe == "" {
		exe = "python"
	}

	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".py") {
			if len(args) > 0 {
				next := cmdargs.New([]string{trimmed}).Append(args...).String()
				return exec.NewContext(ctx, exe, next)
			}

			return exec.NewContext(ctx, exe, trimmed)
		}
	}

	splat := []string{"-c", script}
	return exec.NewContext(ctx, exe, splat...)
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

func RubyScriptContext(ctx context.Context, script string, args ...string) *exec.Cmd {
	exe, _ := exec.Find("ruby", nil)
	if exe == "" {
		exe = "ruby"
	}
	if !strings.ContainsAny(script, "\n\r") {
		trimmed := strings.TrimSpace(script)
		if strings.HasSuffix(trimmed, ".rb") {
			if len(args) > 0 {
				next := cmdargs.New([]string{trimmed}).Append(args...).String()
				return exec.NewContext(ctx, exe, next)
			}
			return exec.NewContext(ctx, exe, trimmed)
		}
	}

	splat := []string{"-e", script}
	return exec.NewContext(ctx, exe, splat...)
}

type ScriptAttributes struct {
	Env    map[string]string
	Dir    string
	Output bool
}
