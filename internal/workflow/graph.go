package workflow

import (
	"errors"

	"github.com/hyprxlabs/xtasks/internal/schema"
)

func flattenTasks(targets []string, tasks schema.Tasks, set []schema.TaskDef) ([]schema.TaskDef, error) {

	for _, target := range targets {
		task, ok := tasks[target]
		if !ok {
			return nil, errors.New("Task not found: " + target)
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

		if len(task.Needs) > 0 {
			neededTasks, err := flattenTasks(task.Needs, tasks, set)
			if err != nil {
				return nil, err
			}
			set = neededTasks
		}
	}

	return set, nil
}

func findCyclicalReferences(tasks []schema.TaskDef) []schema.TaskDef {
	stack := []schema.TaskDef{}
	cycles := []schema.TaskDef{}

	var resolve func(task schema.TaskDef) bool
	resolve = func(task schema.TaskDef) bool {
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
