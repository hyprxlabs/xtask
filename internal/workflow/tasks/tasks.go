package tasks

import (
	"context"
	"strings"
	"time"

	"github.com/hyprxlabs/xtask/internal/errors"
	"github.com/hyprxlabs/xtask/internal/schema"
	"github.com/hyprxlabs/xtask/internal/workflow/statuses"
)

type TaskContext struct {
	Task    Task
	TaskDef schema.TaskDef
	Targets schema.Hosts
	Context context.Context
}

type Task struct {
	Id      string
	Name    string
	Desc    string
	Env     map[string]string
	XEnv    map[string]string
	Run     string
	Uses    string
	Hosts   []string
	Cwd     string
	Timeout time.Duration
	Needs   []string
	Files   []string
	Args    []string
}

type TaskResult struct {
	Err       error
	Status    int
	StartedAt time.Time
	EndedAt   time.Time
	Message   string
	Output    map[string]interface{}
}

func (tr *TaskResult) Start() *TaskResult {
	tr.StartedAt = time.Now().UTC()
	return tr
}

func (tr *TaskResult) End() *TaskResult {
	tr.EndedAt = time.Now().UTC()
	return tr
}

func (tr *TaskResult) Ok() *TaskResult {
	tr.Status = statuses.Ok
	tr.End()
	return tr
}

func (tr *TaskResult) Fail(err error) *TaskResult {
	tr.Err = errors.WithCause(err, tr.Err)
	tr.Status = statuses.Error
	tr.End()
	return tr
}

func (tr *TaskResult) Skip(msg string) *TaskResult {
	tr.Status = statuses.Skipped
	tr.Message = msg
	tr.End()
	return tr
}

func (tr *TaskResult) Cancel(msg string) *TaskResult {
	tr.Status = statuses.Cancelled
	tr.Message = msg
	tr.End()
	return tr
}

func NewTaskResult() *TaskResult {
	return &TaskResult{
		Err:       nil,
		Status:    statuses.None,
		StartedAt: time.Now().UTC(),
		EndedAt:   time.Now().UTC(),
		Message:   "",
		Output:    make(map[string]interface{}),
	}
}

func Run(ctx TaskContext) *TaskResult {
	// Here you would implement the logic to run the task based on the context
	// This is a placeholder for the actual implementation

	if strings.HasPrefix(ctx.Task.Uses, "scp") {
		// If the task uses SCP, we can run the SCP task
		return runSCP(ctx)
	}

	if strings.HasPrefix(ctx.Task.Uses, "docker") {
		// If the task uses Docker, we can run the Docker task
		return runDocker(ctx)
	}

	if strings.HasPrefix(ctx.Task.Uses, "ssh") {
		return runSSH(ctx)
	}

	if ctx.Task.Run != "" {
		// If a command is defined, run it
		return runShell(ctx)
	}

	res := NewTaskResult()
	return res.Fail(errors.NewDetails("Task definition is missing", "task_definition_missing", "The task definition is required to run the task"))
}
