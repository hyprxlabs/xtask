package workflows

import (
	"errors"

	"github.com/hyprxlabs/xtask/types"
)

func flattenTasks(targets []string, tasks types.Tasks, set []types.Task) ([]types.Task, error) {

	for _, target := range targets {
		task, ok := tasks[target]
		if !ok {
			return nil, errors.New("Task not found: " + target)
		}

		if len(task.Needs) > 0 {
			neededTasks, err := flattenTasks(task.Needs, tasks, set)
			if err != nil {
				return nil, err
			}
			set = neededTasks
		}

		added := false
		for _, task2 := range set {
			if task.Id == task2.Id {
				added = true
				break
			}
		}

		if !added {
			set = append(set, task)
		}
	}

	return set, nil
}

func findCyclicalReferences(tasks []types.Task) []types.Task {
	stack := []types.Task{}
	cycles := []types.Task{}

	var resolve func(task types.Task) bool
	resolve = func(task types.Task) bool {
		for _, t := range stack {
			if task.Id == t.Id {
				return false
			}
		}

		stack = append(stack, task)

		if len(task.Needs) > 0 {
			for _, need := range task.Needs {
				for _, nextTask := range tasks {
					if nextTask.Id == need {
						if !resolve(nextTask) {
							return false
						}
					}
				}
			}
		}

		stack = stack[:len(stack)-1]
		return true
	}

	for _, task := range tasks {
		if !resolve(task) {
			cycles = append(cycles, task)
		}
	}

	return cycles
}
