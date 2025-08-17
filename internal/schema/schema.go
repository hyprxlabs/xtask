package schema

import (
	"errors"
	"os"
	"strings"

	"github.com/hyprxlabs/go/env"
	"gopkg.in/yaml.v3"
)

type Hosts map[string]SshTarget
type Tasks map[string]TaskDef

func (t *Tasks) UnmarshalYAML(value *yaml.Node) error {
	if (*t) == nil {
		*t = make(Tasks)
	}

	if value.Kind != yaml.MappingNode {
		return errors.New("tasks must be a mapping of strings to TaskDef")
	}

	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.New("expected string key for task")
		}
		key := keyNode.Value
		valueNode := value.Content[i+1]

		if valueNode.Kind == yaml.ScalarNode {
			if valueNode.Value == "" {
				return errors.New("task definition cannot be empty")
			}

			// If the value is a scalar, we assume it's a command to run
			(*t)[key] = TaskDef{
				Id:     key,
				Name:   &key,
				Env:    make(map[string]string),
				Dotenv: []string{},
				Run:    &valueNode.Value,
			}
			continue
		}

		if valueNode.Kind != yaml.MappingNode {
			return errors.New("expected mapping for task value")
		}

		task := &TaskDef{
			Id: key,
		}
		if err := valueNode.Decode(&task); err != nil {
			return err
		}

		task.Id = key
		if task.Name == nil {
			task.Name = &key
		}

		(*t)[key] = *task
	}

	return nil
}

func (t *Hosts) UnmarshalYAML(value *yaml.Node) error {
	if (*t) == nil {
		*t = make(Hosts)
	}

	if value.Kind == yaml.ScalarNode {
		path := strings.TrimSpace(value.Value)
		if path == "" {
			return errors.New("targets cannot be an empty string")
		}

		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			p2, err := env.Expand(path)
			if err != nil {
				return err
			}
			path = strings.TrimSpace(p2)
			if len(path) > 3 && path[0] == '~' && path[1] == '/' {
				path = env.Get(env.HOME) + path[2:]
			}
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(bytes, t); err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("targets must be a mapping of strings to SshTarget or a string that is a path to a YAML file for targets")
		}
	}

	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content); i += 2 {
			key := value.Content[i]
			if key.Kind != yaml.ScalarNode {
				return errors.New("expected string key for target")
			}
			valueNode := value.Content[i+1]
			if valueNode.Kind == yaml.ScalarNode {
				value := valueNode.Value
				if strings.ContainsRune(value, '@') {
					parts := strings.SplitN(value, "@", 2)
					if len(parts) != 2 {
						return errors.New("invalid SSH target format, expected 'user@host'")
					}
					(*t)[key.Value] = SshTarget{
						Host: parts[1],
						User: &parts[0],
					}
				} else {

					(*t)[key.Value] = SshTarget{
						Host: key.Value,
					}
				}

				continue
			}

			if valueNode.Kind != yaml.MappingNode {
				return errors.New("expected mapping or string for target value")
			}

			var target SshTarget

			if err := value.Content[i+1].Decode(&target); err != nil {
				return err
			}

			(*t)[key.Value] = target
		}
	}

	return errors.New("targets must be a mapping of strings to SshTarget or a string")
}

type Workflow struct {
	Imports        []ImportDef       `yaml:"imports,omitempty"`
	Name           *string           `yaml:"name"`
	Env            map[string]string `yaml:"env,omitempty"`
	Dotenv         []string          `yaml:"dotenv,omitempty"`
	ExpandCommands bool              `yaml:"expand-commands,omitempty"`
	Tasks          Tasks             `yaml:"tasks,omitempty"`
	Hosts          Hosts             `yaml:"hosts,omitempty"`
	Shell          string            `yaml:"shell,omitempty"`
	Path           *Path             `yaml:"path,omitempty"`
}

type Path struct {
	Win   []string `yaml:"win,omitempty"`
	Posix []string `yaml:"posix,omitempty"`
	Macos []string `yaml:"macos,omitempty"`
}

type ImportDef struct {
	Uri       string `yaml:"uri"`
	Optional  bool   `yaml:"optional,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

type TaskDef struct {
	Id      string            `yaml:"id,omitempty"`
	Desc    *string           `yaml:"desc,omitempty"`
	Name    *string           `yaml:"name,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
	Dotenv  []string          `yaml:"dotenv,omitempty"`
	Cwd     *string           `yaml:"cwd,omitempty"`
	Timeout *string           `yaml:"timeout,omitempty"`
	Run     *string           `yaml:"run,omitempty"`
	Uses    *string           `yaml:"uses,omitempty"`
	Args    []string          `yaml:"args,omitempty"`
	Ssh     *string           `yaml:"ssh,omitempty"`
	Needs   []string          `yaml:"needs,omitempty"`
	Scp     []string          `yaml:"scp,omitempty"`
	Targets []string          `yaml:"targets,omitempty"`
	Files   []string          `yaml:"files,omitempty"`
	With    map[string]string `yaml:"with,omitempty"`
}

type SshTarget struct {
	Host     string  `yaml:"host"`
	Port     *int    `yaml:"port,omitempty"`
	User     *string `yaml:"user,omitempty"`
	Identity *string `yaml:"identity,omitempty"`
	// this must point to an env variable that contains the password
	Password *string  `yaml:"password,omitempty"`
	Groups   []string `yaml:"groups,omitempty"`
}
