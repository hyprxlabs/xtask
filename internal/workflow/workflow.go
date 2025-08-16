package workflow

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/hyprxlabs/go/dotenv"
	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtasks/internal/schema"
	"github.com/hyprxlabs/xtasks/internal/workflow/tasks"
	"gopkg.in/yaml.v3"
)

type Params struct {
	Args                []string
	Tasks               []string
	Timeout             time.Duration
	CommandSubstitution bool
	Context             context.Context
}

type WorkflowContext struct {
	Env     map[string]string
	Vars    map[string]interface{}
	Secrets map[string]string
	Targets schema.Targets
	Timeout time.Duration
	Schema  *schema.Workflow
	Tasks   []schema.TaskDef
	Cwd     string
	Context context.Context
}

func Run(tasksFile string, params Params) error {
	// Load the tasks file and parse it

	bytes, err := os.ReadFile(tasksFile)
	if err != nil {
		return err
	}

	plan := &schema.Workflow{}
	if err := yaml.Unmarshal(bytes, plan); err != nil {
		return err
	}

	tasks := []schema.TaskDef{}
	if len(params.Tasks) == 0 {
		params.Tasks = []string{"default"}
	}

	for _, taskName := range params.Tasks {
		next, ok := plan.Tasks[taskName]
		if ok {
			tasks = append(tasks, next)
		} else {
			return errors.New("Task not found: " + taskName)
		}
	}

	ctx := WorkflowContext{
		Env:     make(map[string]string),
		Vars:    make(map[string]interface{}),
		Secrets: make(map[string]string),
		Targets: plan.Targets,
		Context: params.Context,
	}

	doc := dotenv.NewDocument()
	if len(plan.Dotenv) > 0 {
		for _, file := range plan.Dotenv {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			nextDoc, err := dotenv.Parse(string(data))
			if err != nil {
				return err
			}
			doc.Merge(nextDoc)
		}
	}

	if len(plan.Env) > 0 {
		for key, value := range plan.Env {
			doc.Set(key, value)
		}
	}

	ctx.Env = env.All()
	get := func(key string) string {
		if val, ok := ctx.Env[key]; ok {
			return val
		}
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return ""
	}
	set := func(key, value string) error {
		ctx.Env[key] = value
		return nil
	}

	options := &env.ExpandOptions{
		CommandSubstitution: true,
		ShellArgs:           params.Args,
		Get:                 get,
		Set:                 set,
	}

	for _, node := range doc.ToArray() {
		if node.Type == dotenv.VARIABLE_TOKEN {
			key := ""
			value := node.Value
			if node.Key != nil {
				key = *node.Key
			}

			value, err := env.ExpandWithOptions(value, options)
			if err != nil {
				return err
			}
			ctx.Env[key] = value
		}
	}

	ctx.Schema = plan
	ctx.Tasks = tasks
	ctx.Cwd = filepath.Dir(tasksFile)

	os.Chdir(ctx.Cwd)

	return runTasks(ctx)
}

func runTasks(ctx WorkflowContext) error {
	taskDefs := ctx.Tasks
	if len(taskDefs) == 0 {
		return errors.New("no tasks to run")
	}

	set := []tasks.Task{}

	for _, taskDef := range taskDefs {
		name := taskDef.Id
		if taskDef.Name != nil {
			name = *taskDef.Name
		}
		uses := ""
		if taskDef.Uses != nil {
			uses = *taskDef.Uses
		}

		desc := ""
		if taskDef.Desc != nil {
			desc = *taskDef.Desc
		}

		timeout := time.Duration(0)
		if taskDef.Timeout != nil {
			timeout, _ = time.ParseDuration(*taskDef.Timeout)
		}

		cwd := ctx.Cwd
		if taskDef.Cwd != nil {
			cwd = *taskDef.Cwd
		}

		if cwd == "" {
			cwd = ctx.Cwd
		}

		doc := dotenv.NewDocument()

		envMap := make(map[string]string)
		for key, value := range ctx.Env {
			envMap[key] = value
		}

		if len(taskDef.Dotenv) > 0 {
			for _, denv := range taskDef.Dotenv {
				data, err := os.ReadFile(denv)
				if err != nil {
					return errors.New("Failed to read dotenv file: " + err.Error())
				}
				doc, err := dotenv.Parse(string(data))
				if err != nil {
					return errors.New("Failed to parse dotenv file: " + err.Error())
				}

				doc.Merge(doc)
			}
		}

		if len(taskDef.Env) > 0 {
			for key, value := range taskDef.Env {

				doc.Set(key, value)
			}
		}

		if doc.Len() > 0 {

			options := &env.ExpandOptions{
				CommandSubstitution: ctx.Schema.ExpandCommands,
				ShellArgs:           ctx.Schema.Dotenv,
				Get: func(key string) string {
					if val, ok := envMap[key]; ok {
						return val
					}
					if val, ok := os.LookupEnv(key); ok {
						return val
					}
					return ""
				},
				Set: func(key, value string) error {
					envMap[key] = value
					return nil
				},
			}

			for _, node := range doc.ToArray() {
				if node.Type == dotenv.VARIABLE_TOKEN {
					key := ""
					value := node.Value
					if node.Key != nil {
						key = *node.Key
					}

					value, err := env.ExpandWithOptions(value, options)
					if err != nil {
						return errors.New("Failed to expand environment variable: " + err.Error())
					}
					envMap[key] = value
				}
			}
		}

		run := ""
		if taskDef.Run != nil {
			run = *taskDef.Run
		}

		task := &tasks.Task{
			Id:      taskDef.Id,
			Name:    name,
			Desc:    desc,
			Uses:    uses,
			Env:     envMap,
			Run:     run,
			Targets: taskDef.Targets,
			Cwd:     cwd,
			Timeout: timeout,
			Needs:   taskDef.Needs,
			Files:   taskDef.Files,
		}

		set = append(set, *task)
	}

	for _, task := range set {
		var def schema.TaskDef
		for _, t := range taskDefs {
			if t.Id == task.Id {
				def = t
				break
			}
		}

		c := &tasks.TaskContext{
			Task:    task,
			Targets: ctx.Targets,
			Context: ctx.Context,
			TaskDef: def,
		}

		res := tasks.Run(*c)
		if res.Err != nil {
			return res.Err
		}
	}

	return nil
}
