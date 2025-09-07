package workflows

import (
	"context"
	"os"

	"github.com/hyprxlabs/xtask/types"
)

type Workflow struct {
	Name        *string
	App         *string
	Contexts    []string
	Version     *string
	Config      *types.Config
	Env         *types.Env
	Tasks       map[string]types.Task
	Hosts       map[string]types.Host
	Values      map[string]interface{}
	Args        []string
	ContextName string
	Context     context.Context
	cleanupEnv  bool
	cleanupPath bool
}

func NewWorkflow() *Workflow {

	defaultShell := os.Getenv("XTASK_SHELL")
	if len(defaultShell) == 0 {
		defaultShell = "bash"
		if os.Getenv("OS") == "Windows_NT" {
			defaultShell = "powershell"
		}
	}

	defaultContext := os.Getenv("XTASK_CONTEXT")
	if len(defaultContext) == 0 {
		defaultContext = "default"
	}

	return &Workflow{
		Config: &types.Config{
			Substitution: true,
			Dirs: types.Dirs{
				Etc:  "./.xtask/etc",
				Apps: []string{"./.xtask/apps"},
			},
			PrependPaths: []types.PrependPath{},
			Env:          *types.NewEnv(),
			Shell:        defaultShell,
		},
		Env:         types.NewEnv(),
		Values:      map[string]interface{}{},
		ContextName: defaultContext,
		Name:        nil,
		App:         nil,
		Contexts:    []string{"default"},
		Version:     nil,
		Tasks:       map[string]types.Task{},
		Hosts:       map[string]types.Host{},
		Args:        []string{},
		Context:     context.Background(),
		cleanupEnv:  false,
		cleanupPath: false,
	}
}
