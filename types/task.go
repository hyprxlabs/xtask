package types

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Task struct {
	Id        string                 `yaml:"id,omitempty"`
	Desc      *string                `yaml:"desc,omitempty"`
	Help      *string                `yaml:"help,omitempty"`
	Name      *string                `yaml:"name,omitempty"`
	Env       Env                    `yaml:"env,omitempty"`
	Dotenv    []string               `yaml:"dotenv,omitempty"`
	Cwd       *string                `yaml:"cwd,omitempty"`
	Timeout   *string                `yaml:"timeout,omitempty"`
	Run       *string                `yaml:"run,omitempty"`
	Uses      *string                `yaml:"uses,omitempty"`
	Args      []string               `yaml:"args,omitempty"`
	Needs     []string               `yaml:"needs,omitempty"`
	Hosts     []string               `yaml:"hosts,omitempty"`
	With      map[string]interface{} `yaml:"with,omitempty"`
	Predicate *string                `yaml:"if,omitempty"`
}

type Tasks map[string]Task

func (t *Tasks) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	if t == nil {
		t = &Tasks{}
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return errors.New("task key must be a string")
		}

		key := keyNode.Value

		var task Task
		if valueNode.Kind == yaml.MappingNode {
			if err := valueNode.Decode(&task); err != nil {
				return err
			}
			task.Id = key
			(*t)[key] = task
			continue
		}

		if valueNode.Kind == yaml.ScalarNode {
			task = Task{
				Id:     key,
				Run:    &valueNode.Value,
				Env:    Env{},
				With:   map[string]interface{}{},
				Dotenv: []string{},
				Args:   []string{},
				Needs:  []string{},
				Hosts:  []string{},
			}

			(*t)[key] = task
			continue
		}
	}

	return nil
}

type SharedTask struct {
	Id     string  `yaml:"id,omitempty"`
	Desc   *string `yaml:"desc,omitempty"`
	Help   *string `yaml:"help,omitempty"`
	Run    *string `yaml:"run,omitempty"`
	Uses   *string `yaml:"uses,omitempty"`
	Inputs []Input `yaml:"inputs,omitempty"`
}

type Output struct {
	Id      string  `yaml:"id,omitempty"`
	Desc    *string `yaml:"desc,omitempty"`
	Default *string `yaml:"default,omitempty"`
	Type    *string `yaml:"type,omitempty"`
}

type Input struct {
	Id       string  `yaml:"id,omitempty"`
	Name     *string `yaml:"name,omitempty"`
	Desc     *string `yaml:"desc,omitempty"`
	Help     *string `yaml:"help,omitempty"`
	Default  *string `yaml:"default,omitempty"`
	Type     *string `yaml:"type,omitempty"`
	Required *bool   `yaml:"required,omitempty"`
}
