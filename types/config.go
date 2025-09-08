package types

import (
	"errors"
	"os"
	"strings"

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
	Etc     string   `yaml:"etc,omitempty" mapstructure:"etc,omitempty"`
	Apps    []string `yaml:"apps,omitempty" mapstructure:"apps,omitempty"`
	Scripts string   `yaml:"scripts,omitempty" mapstructure:"scripts,omitempty"`
	Bin     string   `yaml:"bin,omitempty" mapstructure:"bin,omitempty"`
}

func (d *Dirs) UnmarshalYAML(node *yaml.Node) error {

	if d == nil {
		d = &Dirs{}
	}

	dotdir := os.Getenv("XTASK_DOT_DIR")
	if len(dotdir) == 0 {
		dotdir = "./.xtask"
	}

	if d.Etc == "" {
		etc := os.Getenv("XTASK_ETC_DIR")
		if len(etc) == 0 {
			etc = dotdir + "/etc"
		}
		d.Etc = etc
	}

	if d.Apps == nil || len(d.Apps) == 0 {
		apps := os.Getenv("XTASK_APPS_DIRS")
		if len(apps) == 0 {
			apps = dotdir + "/apps"
		}
		d.Apps = strings.Split(apps, string(os.PathListSeparator))
	}

	if d.Scripts == "" {
		scripts := os.Getenv("XTASK_SCRIPTS_DIR")
		if len(scripts) == 0 {
			scripts = dotdir + "/scripts"
		}
		d.Scripts = scripts
	}

	if d.Bin == "" {
		bin := os.Getenv("XTASK_BIN_DIR")
		if len(bin) == 0 {
			bin = dotdir + "/bin"
		}
		d.Bin = bin
	}

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
		case "scripts":
			if valueNode.Kind != yaml.ScalarNode {
				return errors.New("expected string value for dirs.scripts")
			}
			d.Scripts = valueNode.Value
		case "bin":
			if valueNode.Kind != yaml.ScalarNode {
				return errors.New("expected string value for dirs.bin")
			}
			d.Bin = valueNode.Value
		default:
			return errors.New("unknown key in dirs: " + key)
		}
	}
	return nil
}
