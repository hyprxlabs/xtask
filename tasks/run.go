package tasks

import (
	"net/url"
	"strings"

	"github.com/hyprxlabs/xtask/errors"
)

func Run(ctx TaskContext) *TaskResult {
	// Here you would implement the logic to run the task based on the context
	// This is a placeholder for the actual

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
	case "bash", "sh", "zsh", "powershell", "pwsh", "cmd", "python", "ruby":
		return runShell(ctx)
	default:
		// Unsupported task type
		res := NewTaskResult()
		return res.Fail(errors.NewDetails("Unsupported task type: "+uses, "unsupported_task_type", "The task type is not supported"))
	}
}
