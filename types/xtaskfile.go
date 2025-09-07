package types

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type XTaskfile struct {
	Path      string                 `yaml:"-"`
	Imports   Imports                `yaml:"imports,omitempty"`
	Name      *string                `yaml:"name"`
	App       *string                `yaml:"app,omitempty"`
	Contexts  []string               `yaml:"context,omitempty"`
	Version   *string                `yaml:"version,omitempty"`
	Config    *Config                `yaml:"config,omitempty"`
	Env       *Env                   `yaml:"env,omitempty"`
	Secrets   *Env                   `yaml:"secrets,omitempty"`
	Dotenv    []string               `yaml:"dotenv,omitempty"`
	Tasks     *Tasks                 `yaml:"tasks,omitempty"`
	HostsNode *HostsNode             `yaml:"hosts,omitempty"`
	Values    map[string]interface{} `yaml:"values,omitempty"`
}

func NewXTaskfile() *XTaskfile {
	defaultContext := os.Getenv("XTASK_CONTEXT")
	if len(defaultContext) == 0 {
		defaultContext = "default"
	}

	defaultShell := os.Getenv("XTASK_SHELL")
	if len(defaultShell) == 0 {
		defaultShell = "bash"
		if os.Getenv("OS") == "Windows_NT" {
			defaultShell = "powershell"
		}
	}

	return &XTaskfile{
		Imports:  Imports{},
		Name:     nil,
		App:      nil,
		Contexts: []string{defaultContext},
		Version:  nil,
		Config: &Config{
			Substitution: true,
			Dirs: Dirs{
				Etc:  "./.xtask/etc",
				Apps: []string{"./.xtask/apps"},
			},
			PrependPaths: []PrependPath{},
			Env:          Env{},
			Shell:        defaultShell,
		},
		Env:       NewEnv(),
		Secrets:   NewEnv(),
		Dotenv:    []string{},
		Tasks:     &Tasks{},
		HostsNode: &HostsNode{Hosts: Hosts{}, Imports: []string{}},
		Values:    map[string]interface{}{},
	}
}

func (x *XTaskfile) DecodeYAMLFile(path string) error {
	if x == nil {
		x = NewXTaskfile()
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		path = absPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, x); err != nil {
		return err
	}

	x.Path = path
	return nil
}
