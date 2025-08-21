package schema

import (
	"errors"
	"os"
	"strings"

	"github.com/hyprxlabs/go/env"
	"gopkg.in/yaml.v3"
)

type Hosts map[string]SshHost
type Tasks map[string]TaskDef

type PrependPaths []PrependedPath

func (p *PrependPaths) UnmarshalYAML(value *yaml.Node) error {
	if (*p) == nil {
		*p = make(PrependPaths, 0)
	}

	if value.Kind != yaml.SequenceNode {
		return errors.New("prepend-paths must be a sequence of strings")
	}

	for _, item := range value.Content {
		if item.Kind != yaml.MappingNode && item.Kind != yaml.ScalarNode {
			return errors.New("prepend-paths must be a sequence of strings or mappings")
		}

		if item.Kind == yaml.ScalarNode {
			path := &PrependedPath{
				Path: strings.TrimSpace(item.Value),
			}

			*p = append(*p, *path)
		}

		if item.Kind == yaml.MappingNode {
			path := &PrependedPath{}
			for i := 0; i < len(item.Content); i += 2 {
				keyNode := item.Content[i]
				if keyNode.Kind != yaml.ScalarNode {
					return errors.New("expected string key for prepend-paths")
				}

				key := keyNode.Value
				valueNode := item.Content[i+1]
				if valueNode.Kind != yaml.ScalarNode {
					return errors.New("expected string value for prepend-paths")
				}

				value := strings.TrimSpace(valueNode.Value)
				switch key {
				case "path":
					if path.Path != "" {
						return errors.New("path key already set for prepend-paths")
					}
					if value == "" {
						return errors.New("path cannot be empty for prepend-paths")
					}
					path.Path = value
				case "os", "targets":
					if path.Os != "" {
						return errors.New("os/targets key already set for prepend-paths")
					}
					if value == "" {
						return errors.New("os/targets cannot be empty for prepend-paths")
					}
				case "win":
					fallthrough
				case "windows":
					if path.Os != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.Os = "windows"
					path.Path = value
				case "linux":
					fallthrough
				case "unix":
					fallthrough
				case "posix":
					if path.Os != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.Os = "posix"
					path.Path = value
				case "macos":
					fallthrough
				case "mac":
					fallthrough
				case "osx":
					fallthrough
				case "darwin":
					if path.Os != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.Os = "darwin"
					path.Path = value
				default:
					return errors.New("unknown key in prepend-paths: " + key)
				}
			}
		}
	}

	return nil
}

func (t *Tasks) UnmarshalYAML(value *yaml.Node) error {
	if (*t) == nil {
		*t = make(Tasks)
	}

	if value.Kind != yaml.MappingNode && value.Kind != yaml.ScalarNode {
		return errors.New("tasks must be a mapping of strings to TaskDef")
	}

	if value.Kind == yaml.ScalarNode {
		path := strings.TrimSpace(value.Value)
		if path == "" {
			return errors.New("tasks cannot be an empty string")
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
			return errors.New("tasks must be a mapping of strings to TaskDef or a string that is a path to a YAML file for tasks")
		}
	}

	if value.Kind == yaml.MappingNode {
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
	}

	return nil
}

func (t *Hosts) UnmarshalYAML(value *yaml.Node) error {
	if (*t) == nil {
		*t = make(Hosts)
	}

	if value.Kind != yaml.ScalarNode && value.Kind != yaml.MappingNode {
		return errors.New("hosts must be a mapping or a string and not " + getNodeType(value))
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
					(*t)[key.Value] = SshHost{
						Host: parts[1],
						User: &parts[0],
					}
				} else {

					(*t)[key.Value] = SshHost{
						Host: key.Value,
					}
				}

				continue
			}

			if valueNode.Kind != yaml.MappingNode {
				return errors.New("expected mapping or string for host value")
			}

			var target SshHost

			if err := value.Content[i+1].Decode(&target); err != nil {
				return err
			}

			(*t)[key.Value] = target
		}
	}

	return nil
}

type Config struct {
	// The path adds additional directories to the PATH environment variable
	// This is useful for adding directories where executables are located before
	// any tasks are run or env files are loaded.
	PrependPaths PrependPaths `yaml:"prepend-paths,omitempty"`
	// The env section in the config loads environment variables
	// into the process before any other environment variables, dotenv files,
	// are loaded or tasks are run.
	Env map[string]string `yaml:"env,omitempty"`
	// The directories to search in for delegation running a task
	// by calling `xtask drun [dir] <task>`. The drun command will
	// use the value for [dir] and look in each of the directories
	// to find match for the combined path of [delegation_dir][dir]/xtaskfile.yaml
	DelegationDirs      []string `yaml:"delegation-dirs,omitempty"`
	Shell               string   `yaml:"shell,omitempty"`
	CommandSubstitution bool     `yaml:"command-substitution,omitempty"`
}

type Workflow struct {
	Imports []ImportDef       `yaml:"imports,omitempty"`
	Name    *string           `yaml:"name"`
	Config  *Config           `yaml:"config,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
	Dotenv  []string          `yaml:"dotenv,omitempty"`
	Tasks   Tasks             `yaml:"tasks,omitempty"`
	Hosts   Hosts             `yaml:"hosts,omitempty"`
	Values  []string          `yaml:"values,omitempty"`
}

type PrependedPath struct {
	Path string `yaml:"path"`
	Os   string `yaml:"targets,omitempty"`
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
	Hosts   []string          `yaml:"hosts,omitempty"`
	Files   []string          `yaml:"files,omitempty"`
	With    map[string]string `yaml:"with,omitempty"`
}

type SshHost struct {
	Host     string  `yaml:"host"`
	Port     *int    `yaml:"port,omitempty"`
	User     *string `yaml:"user,omitempty"`
	Identity *string `yaml:"identity,omitempty"`
	// this must point to an env variable that contains the password
	Password *string                `yaml:"password,omitempty"`
	Groups   []string               `yaml:"groups,omitempty"`
	Meta     map[string]interface{} `yaml:"meta,omitempty"`
	Os       *OsInfo                `yaml:"os,omitempty"`
}

type OsInfo struct {
	Platform     string `yaml:"platform"`
	Arch         string `yaml:"arch"`
	Variant      string `yaml:"variant,omitempty"`
	Family       string `yaml:"family,omitempty"`
	Codename     string `yaml:"codename,omitempty"`
	Version      string `yaml:"version,omitempty"`
	BuildVersion string `yaml:"build_version,omitempty"`
}

func getNodeType(n *yaml.Node) string {
	switch n.Kind {
	case yaml.ScalarNode:
		return "scalar"
	case yaml.MappingNode:
		return "mapping"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.AliasNode:
		return "alias"
	case yaml.DocumentNode:
		return "document"
	default:
		return "unknown"
	}
}
