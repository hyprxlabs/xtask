package types

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type PrependPath struct {
	Path string `yaml:"path"`
	OS   string `yaml:"os,omitempty"`
}

type PrependPaths []PrependPath

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
			path := &PrependPath{
				Path: strings.TrimSpace(item.Value),
			}

			*p = append(*p, *path)
		}

		if item.Kind == yaml.MappingNode {
			path := &PrependPath{}
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
					if path.OS != "" {
						return errors.New("os/targets key already set for prepend-paths")
					}
					if value == "" {
						return errors.New("os/targets cannot be empty for prepend-paths")
					}
				case "win":
					fallthrough
				case "windows":
					if path.OS != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.OS = "windows"
					path.Path = value
				case "linux":
					fallthrough
				case "unix":
					fallthrough
				case "posix":
					if path.OS != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.OS = "posix"
					path.Path = value
				case "macos":
					fallthrough
				case "mac":
					fallthrough
				case "osx":
					fallthrough
				case "darwin":
					if path.OS != "" {
						return errors.New("os key already set for prepend-paths")
					}
					path.OS = "darwin"
					path.Path = value
				default:
					return errors.New("unknown key in prepend-paths: " + key)
				}
			}
		}
	}

	return nil
}
