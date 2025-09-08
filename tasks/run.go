package tasks

import (
	"net/url"
	"strings"

	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/go/exec"
	"github.com/hyprxlabs/xtask/errors"
	"github.com/hyprxlabs/xtask/types"
)

type taskEnvLike struct {
	Env *types.Env
}

func (t *taskEnvLike) Get(key string) string {
	if t.Env == nil {
		return ""
	}
	return t.Env.GetString(key)
}

func (t *taskEnvLike) Expand(s string) (string, error) {
	if t.Env == nil {
		return s, nil
	}
	opts := env.ExpandOptions{
		Get: t.Env.GetString,
		Set: func(key, value string) error {
			t.Env.Set(key, value)
			return nil
		},
		Keys:                t.Env.Keys(),
		CommandSubstitution: true,
	}

	return env.ExpandWithOptions(s, &opts)
}

func (t *taskEnvLike) Set(key, value string) {
	if t.Env == nil {
		return
	}
	t.Env.Set(key, value)
}

func (t *taskEnvLike) SplitPath() []string {
	if t.Env == nil {
		return []string{}
	}
	return t.Env.SplitPath()
}

func Run(ctx TaskContext) *TaskResult {

	// Set custom env like for exec package
	// so that the task env is used for finding executables
	oldEnv := exec.GetEnvLike()
	defer exec.SetEnvLike(oldEnv)
	envLike := &taskEnvLike{Env: &ctx.Data.Env}
	exec.SetEnvLike(envLike)

	uses := ctx.Data.Uses
	if strings.Contains(uses, "://") {
		uri, err := url.Parse(uses)

		if err != nil {
			res := NewTaskResult()
			return res.Fail(errors.New("Invalid template URI: " + err.Error()))
		}

		uses = uri.Scheme
	}

	switch uses {
	case "tmpl":
		return runTpl(ctx)
	case "scp":
		return runSCP(ctx)
	case "ssh":
		return runSSH(ctx)
	case "docker":
		return runDocker(ctx)
	case "bash", "sh", "zsh", "powershell", "pwsh", "cmd", "python", "ruby", "deno", "node", "bun":
		return runShell(ctx)
	default:
		// Unsupported task type
		res := NewTaskResult()
		return res.Fail(errors.NewDetails("Unsupported task type: "+uses, "unsupported_task_type", "The task type is not supported"))
	}
}
