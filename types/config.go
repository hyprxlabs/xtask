package types

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// The path adds additional directories to the PATH environment variable
	// This is useful for adding directories where executables are located before
	// any tasks are run or env files are loaded.
	PrependPaths PrependPaths `yaml:"prepend-paths,omitempty" mapstructure:"prepend-paths,omitempty"`
	// The env section in the config loads environment variables
	// into the process before any other environment variables, dotenv files,
	// are loaded or tasks are run.
	Env Env `yaml:"env,omitempty" mapstructure:"env,omitempty"`
	// The directories to search in for delegation running a task
	// by calling `xtask drun [dir] <task>`. The drun command will
	// use the value for [dir] and look in each of the directories
	// to find match for the combined path of [delegation_dir][dir]/xtaskfile.yaml
	Dirs         Dirs    `yaml:"dirs,omitempty" mapstructure:"dirs,omitempty"`
	Shell        string  `yaml:"shell,omitempty" mapstructure:"shell,omitempty"`
	Substitution bool    `yaml:"substitution,omitempty" mapstructure:"substitution,omitempty"`
	Context      *string `yaml:"context,omitempty" mapstructure:"context,omitempty"`
}

type Dirs struct {
	Etc  string   `yaml:"etc,omitempty" mapstructure:"etc,omitempty"`
	Apps []string `yaml:"apps,omitempty" mapstructure:"apps,omitempty"`
	Walk []string `yaml:"walk,omitempty" mapstructure:"walk,omitempty"`
}

func (d *Dirs) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("dirs must be a mapping")
	}
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return errors.New("expected string key for dirs")
		}
		key := keyNode.Value

		switch key {
		case "etc":
			if valueNode.Kind != yaml.ScalarNode {
				return errors.New("expected string value for dirs.etc")
			}
			d.Etc = valueNode.Value
		case "apps":
			if valueNode.Kind == yaml.SequenceNode {
				var apps []string
				if err := valueNode.Decode(&apps); err != nil {
					return err
				}
				d.Apps = apps
			} else if valueNode.Kind == yaml.ScalarNode {
				d.Apps = []string{valueNode.Value}
			} else {
				return errors.New("expected string or sequence of strings for dirs.apps")
			}
		case "walk":
			if valueNode.Kind == yaml.SequenceNode {
				var walk []string
				if err := valueNode.Decode(&walk); err != nil {
					return err
				}
				d.Walk = walk
			} else if valueNode.Kind == yaml.ScalarNode {
				d.Walk = []string{valueNode.Value}
			} else {
				return errors.New("expected string or sequence of strings for dirs.walk")
			}
		default:
			return errors.New("unknown key in dirs: " + key)
		}
	}
	return nil
}
