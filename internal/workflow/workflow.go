package workflow

import (
	"bufio"
	"context"
	"errors"
	"os"
	ex "os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hyprxlabs/go/dotenv"
	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/go/exec"
	"github.com/hyprxlabs/xtask/internal/schema"
	"github.com/hyprxlabs/xtask/internal/workflow/tasks"
	"gopkg.in/yaml.v3"
)

type Params struct {
	File                string
	Args                []string
	Tasks               []string
	Timeout             time.Duration
	CommandSubstitution bool
	Context             context.Context
	Command             string
	Dotenv              []string
	Env                 map[string]string
}

type WorkflowContext struct {
	Env     map[string]string
	Vars    map[string]interface{}
	Secrets map[string]string
	Targets schema.Hosts
	Timeout time.Duration
	Schema  *schema.Workflow
	Tasks   []schema.TaskDef
	Cwd     string
	Context context.Context
	Args    []string
	Shell   string
}

func Run(params Params) error {
	// Load the tasks file and parse it

	cleanUpEnv := false
	cleanUpPath := false
	tasksFile := params.File
	if !filepath.IsAbs(tasksFile) {
		absolute, err := filepath.Abs(tasksFile)
		if err != nil {
			return err
		}
		tasksFile = absolute
		params.File = tasksFile
	}

	bytes, err := os.ReadFile(tasksFile)
	if err != nil {
		return err
	}

	plan := &schema.Workflow{}
	if err := yaml.Unmarshal(bytes, plan); err != nil {
		return err
	}

	taskSet := []schema.TaskDef{}
	if len(plan.Tasks) > 0 {
		for _, task := range plan.Tasks {
			taskSet = append(taskSet, task)
		}
	}

	switch params.Command {
	case "ls":
		fallthrough
	case "list":
		max := 0
		for _, task := range taskSet {
			if len(task.Id) > max {
				max = len(task.Id)
			}
		}

		for _, task := range taskSet {
			id := task.Id
			desc := ""
			if task.Desc != nil {
				desc = *task.Desc
			}

			if len(desc) > 0 {
				// Pad the ID to align the output
				id = id + " "
				for len(id) < max {
					id += " "
				}

				println("\x1b[34m" + id + "- " + desc + "\x1b[39m")
			} else {
				println("\x1b[34m" + id + "\x1b[39m")
			}

		}

		return nil
	case "exec":
		execCommand(plan, params)
		return nil
	}

	cycles := findCyclicalReferences(taskSet)
	if len(cycles) > 0 {
		println("Cyclical references found in tasks:")
		for _, cycle := range cycles {
			println(" - " + cycle.Id)
		}

		os.Exit(1)
	}

	tasks := []schema.TaskDef{}
	if len(params.Tasks) == 0 {
		params.Tasks = []string{"default"}
	}

	tasks, err = flattenTasks(params.Tasks, plan.Tasks, tasks)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return errors.New("no tasks found to run")
	}

	ctx := WorkflowContext{
		Env:     make(map[string]string),
		Vars:    make(map[string]interface{}),
		Secrets: make(map[string]string),
		Targets: plan.Hosts,
		Context: params.Context,
		Args:    params.Args,
		Shell:   plan.Shell,
	}

	// handle the special case of windows causing max pain when using bash
	// by injecting a shim of bash.exe that maps to WSL and c:\\Windows\\System32\\bash.exe
	// is generally higher in the PATH prority than the git and the path
	// to c:\\Program Files\\Git\\bin gets appended at the end of the PATH
	// everytime git is installed or updated
	if runtime.GOOS == "windows" {
		if ctx.Shell == "bash" || ctx.Shell == "sh" {
			if ctx.Shell == "bash" {
				path, _ := exec.Which("bash")
				if path != "" && strings.EqualFold("C:\\Windows\\System32\\bash.exe", path) {
					gitBin := "C:\\Program Files\\Git\\bin"
					if _, err := os.Stat(gitBin); err == nil {
						env.PrependPath(gitBin)
					}

					// handle the case where OpenSSH is installed
					// and needs to be preended above git to avoid using
					// the git version of ssh.exe.
					openSSH := "C:\\Program Files\\OpenSSH"
					if _, err := os.Stat(openSSH); err == nil {
						env.PrependPath(openSSH)
					}
				}
			}

			if ctx.Shell == "sh" {
				gitBin := "C:\\Program Files\\Git\\bin"
				if _, err := os.Stat(gitBin); err == nil {
					env.PrependPath(gitBin)
				}

				// handle the case where OpenSSH is installed
				// and needs to be preended above git to avoid using
				// the git version of ssh.exe.
				openSSH := "C:\\Program Files\\OpenSSH"
				if _, err := os.Stat(openSSH); err == nil {
					env.PrependPath(openSSH)
				}
			}
		}
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

	// params must override the plan's dotenv and env
	if len(params.Dotenv) > 0 {
		for _, file := range params.Dotenv {
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

	if len(params.Env) > 0 {
		for key, value := range params.Env {
			doc.Set(key, value)
		}
	}

	ctx.Env = env.All()
	envPath := env.Get(env.PATH)
	ctx.Env = normalizeEnvMap(ctx.Env)
	ctx.Env["XTASK_FILE"] = tasksFile
	ctx.Env["XTASK_DIR"] = filepath.Dir(tasksFile)
	ctx.Env["XTASK_SHELL"] = ctx.Shell

	if _, ok := ctx.Env["XTASK_ENV"]; !ok {
		f, err := os.CreateTemp("", "xtask-env-")
		if err != nil {
			return err
		}
		f.Write([]byte{})
		f.Close()
		ctx.Env["XTASK_ENV"] = f.Name()
		cleanUpEnv = true
	}

	if _, ok := ctx.Env["XTASK_PATH"]; !ok {
		f, err := os.CreateTemp("", "xtask-path-")
		if err != nil {
			return err
		}
		f.Write([]byte{})
		f.Close()

		ctx.Env["XTASK_PATH"] = f.Name()
		cleanUpPath = true
	}

	defer func() {
		if cleanUpEnv {
			envFile := ctx.Env["XTASK_ENV"]
			if envFile != "" {
				os.Remove(envFile)
			}
		}
	}()

	defer func() {
		if cleanUpPath {
			pathFile := ctx.Env["XTASK_PATH"]
			if pathFile != "" {
				os.Remove(pathFile)
			}
		}
	}()

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

			if strings.HasPrefix("XTASK", key) {
				continue
			}

			value, err := env.ExpandWithOptions(value, options)
			if err != nil {
				return err
			}
			ctx.Env[key] = value
		}
	}

	ctx.Env[env.PATH] = envPath

	ctx.Schema = plan
	ctx.Tasks = tasks
	ctx.Cwd = filepath.Dir(tasksFile)

	os.Chdir(ctx.Cwd)

	err = runTasks(ctx)
	return err
}

func runTasks(ctx WorkflowContext) error {
	taskDefs := ctx.Tasks
	if len(taskDefs) == 0 {
		return errors.New("no tasks to run")
	}

	set := []tasks.Task{}

	// capture path
	envPath := env.Get(env.PATH)

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

		if strings.ContainsRune(cwd, '$') {
			cwd, _ = env.Expand(cwd)
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
				doc2, err := dotenv.Parse(string(data))
				if err != nil {
					return errors.New("Failed to parse dotenv file: " + err.Error())
				}

				doc.Merge(doc2)
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

					if strings.HasPrefix("XTASK", key) {
						continue
					}

					value, err := env.ExpandWithOptions(value, options)
					if err != nil {
						return errors.New("Failed to expand environment variable: " + err.Error())
					}
					envMap[key] = value
				}
			}
		}

		envMap[env.PATH] = env.Get(env.PATH)

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
			Args:    ctx.Args,
		}

		if task.Uses == "" {
			task.Uses = ctx.Shell
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

		envFile := ctx.Env["XTASK_ENV"]
		if len(envFile) > 0 {
			canOpen := true
			if _, err := os.Stat(envFile); err != nil {
				canOpen = false
			}

			if canOpen {
				bytes, err := os.ReadFile(envFile)
				if err != nil {
					return errors.New("Failed to read XTASK_ENV file: " + err.Error())
				}

				if len(bytes) > 0 {
					doc2, err := dotenv.Parse(string(bytes))
					if err != nil {
						return errors.New("Failed to parse XTASK_ENV file: " + err.Error())
					}

					for _, node := range doc2.ToArray() {
						if node.Type == dotenv.VARIABLE_TOKEN {
							key := ""
							value := node.Value
							if node.Key != nil {
								key = *node.Key
							}

							if strings.HasPrefix("XTASK", key) {
								continue
							}

							value, err := env.ExpandWithOptions(value, &env.ExpandOptions{
								CommandSubstitution: true,
								ShellArgs:           ctx.Args,
								Get: func(key string) string {
									if val, ok := task.Env[key]; ok {
										return val
									}
									if val, ok := os.LookupEnv(key); ok {
										return val
									}
									return ""
								},
								Set: func(key, value string) error {
									task.Env[key] = value
									return nil
								},
							})
							if err != nil {
								return errors.New("Failed to expand environment variable: " + err.Error())
							}
							task.Env[key] = value
						}
					}
				}
			}
		}

		pathFile := ctx.Env["XTASK_PATH"]
		// reset path
		env.SetPath(envPath)
		if len(pathFile) > 0 {

			canOpen := true
			if _, err := os.Stat(pathFile); err != nil {
				canOpen = false
			}

			if canOpen {
				bytes, err := os.ReadFile(pathFile)
				if err != nil {
					return errors.New("Failed to read XTASK_PATH file: " + err.Error())
				}

				if len(bytes) > 0 {
					content := string(bytes)
					scanner := bufio.NewScanner(strings.NewReader(content))
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if len(line) > 0 {
							if _, err := os.Stat(line); err == nil {
								// LAST IN SHOULD BE FIRST IN PATH
								env.PrependPath(line)
							}
						}
					}
				}
			}
		}

		task.Env[env.PATH] = env.Get(env.PATH)

		c := &tasks.TaskContext{
			Task:    task,
			Targets: ctx.Targets,
			Context: ctx.Context,
			TaskDef: def,
		}

		println("\x1b[1m" + task.Name + "\x1b[22m")
		res := tasks.Run(*c)
		if res.Err != nil {
			return res.Err
		}
	}

	return nil
}

func execCommand(plan *schema.Workflow, params Params) {

	doc := dotenv.NewDocument()
	if len(plan.Dotenv) > 0 {
		for _, file := range plan.Dotenv {
			data, err := os.ReadFile(file)
			if err != nil {
				println("Failed to read dotenv file:", err)
				os.Exit(1)
			}
			nextDoc, err := dotenv.Parse(string(data))
			if err != nil {
				println("Failed to parse dotenv file:", err)
				os.Exit(1)
			}
			doc.Merge(nextDoc)
		}
	}

	if len(plan.Env) > 0 {
		for key, value := range plan.Env {
			doc.Set(key, value)
		}
	}

	// params must override the plan's dotenv and env
	if len(params.Dotenv) > 0 {
		for _, file := range params.Dotenv {
			data, err := os.ReadFile(file)
			if err != nil {
				println("Failed to read dotenv file:", err)
				os.Exit(1)
			}
			nextDoc, err := dotenv.Parse(string(data))
			if err != nil {
				println("Failed to parse dotenv file:", err)
				os.Exit(1)
			}
			doc.Merge(nextDoc)
		}
	}

	if len(params.Env) > 0 {
		for key, value := range params.Env {
			doc.Set(key, value)
		}
	}

	envMap := env.All()

	envMap = normalizeEnvMap(envMap)
	options := &env.ExpandOptions{
		CommandSubstitution: params.CommandSubstitution,
		ShellArgs:           params.Args,
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
				println("Failed to expand environment variable:", err)
				os.Exit(1)
			}
			envMap[key] = value
		}
	}

	cmd := &exec.Cmd{
		Cmd: ex.CommandContext(params.Context, params.Args[0], params.Args[1:]...),
	}

	cmd.Dir = filepath.Dir(params.File)
	cmd.WithEnvMap(envMap)

	o, err := cmd.Run()
	if err != nil || o.Code != 0 {
		if o.Code != 0 {
			os.Exit(o.Code)
		} else {
			os.Exit(1)
		}
	}

	os.Exit(o.Code)
}

func normalizeEnvMap(envMap map[string]string) map[string]string {
	if runtime.GOOS == "windows" {
		envMap["HOME"] = env.Get(env.HOME)
		envMap["USER"] = env.Get(env.USER)
		envMap["HOSTNAME"] = env.Get(env.HOSTNAME)
		envMap["OSTYPE"] = "windows"
	}

	envMap["OS_PLATFORM"] = runtime.GOOS
	envMap["OS_ARCH"] = runtime.GOARCH

	if !env.Has(env.HOME_CONFIG) {
		dir, _ := os.UserConfigDir()
		if dir != "" {
			envMap[env.HOME_CONFIG] = dir
		}
	}

	if !env.Has(env.HOME_CACHE) {
		dir, _ := os.UserCacheDir()
		if dir != "" {
			envMap[env.HOME_CACHE] = dir
		}
	}

	if !env.Has(env.HOME_DATA) {
		switch runtime.GOOS {
		case "windows":
			localAppData := env.Get("LOCALAPPDATA")
			if localAppData != "" {
				envMap[env.HOME_DATA] = localAppData
			} else {
				envMap[env.HOME_DATA] = filepath.Join(env.Get(env.HOME), "AppData", "Local")
			}
		case "darwin":
			envMap[env.HOME_DATA] = filepath.Join(env.Get(env.HOME), "Library", "Application Support")
		default:
			envMap[env.HOME_DATA] = filepath.Join(env.Get(env.HOME), ".local", "share")
		}
	}

	return envMap
}
