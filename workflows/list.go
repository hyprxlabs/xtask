package workflows

import (
	"github.com/hyprxlabs/xtask/types"
)

func (wf *Workflow) List() []types.Task {
	if wf == nil {
		return []types.Task{}
	}

	tasks := []types.Task{}
	for _, task := range wf.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}
