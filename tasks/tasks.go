package tasks

import (
	"context"
	"time"

	"github.com/hyprxlabs/xtask/types"
)

type TaskData struct {
	Id      string
	Help    string
	Desc    string
	Env     types.Env
	Run     string
	Uses    string
	Hosts   types.Hosts
	Cwd     string
	Timeout time.Duration
	Needs   []string
	With    map[string]interface{}
}

type TaskContext struct {
	Task        types.Task
	Data        TaskData
	Context     context.Context
	Args        []string
	ContextName string
}
