package types

import (
	"errors"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Host struct {
	Host     string  `yaml:"host"`
	Port     *int    `yaml:"port,omitempty"`
	User     *string `yaml:"user,omitempty"`
	Identity *string `yaml:"identity,omitempty"`
	// this must point to an env variable that contains the password
	Password *string                `yaml:"password,omitempty"`
	Groups   []string               `yaml:"groups,omitempty"`
	Meta     map[string]interface{} `yaml:"meta,omitempty"`
	OS       *OS                    `yaml:"os,omitempty"`
	Defaults string                 `yaml:"defaults,omitempty"`
}

type HostsNode struct {
	Hosts   Hosts    `yaml:"hosts"`
	Imports []string `yaml:"imports,omitempty"`
}

type Hosts map[string]Host

func (hn *HostsNode) UnmarshalYAML(node *yaml.Node) error {

	if hn == nil {
		hn = &HostsNode{}
	}

	if hn.Imports == nil {
		hn.Imports = []string{}
	}

	if hn.Hosts == nil {
		hn.Hosts = Hosts{}
	}

	if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			if item.Kind == yaml.ScalarNode {
				hn.Imports = append(hn.Imports, item.Value)
				continue
			}

			if item.Kind != yaml.MappingNode {
				continue
			}

			var host Host
			if err := item.Decode(&host); err != nil {
				return err
			}

			if host.Host == "" {
				return errors.New("host entry missing host field")
			}

			hn.Hosts[host.Host] = host
		}

		return nil
	}

	if node.Kind != yaml.MappingNode {
		return errors.New("invalid hosts entry")
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return errors.New("host key must be a string")
		}

		key := keyNode.Value

		if valNode.Kind == yaml.ScalarNode {
			next := &Host{}
			hostname := valNode.Value
			if strings.ContainsRune(valNode.Value, '@') {
				parts := strings.SplitN(valNode.Value, "@", 2)
				next.User = &parts[0]
				hostname = parts[1]
			}

			if strings.ContainsRune(hostname, ':') {
				parts := strings.SplitN(hostname, ":", 2)
				hostname = parts[0]
				if len(parts[1]) > 0 {
					port, err := strconv.Atoi(parts[1])
					if err != nil {
						return err
					}
					next.Port = &port
				}
			}

			next.Host = hostname
			hn.Hosts[key] = *next
			continue
		}

		if valNode.Kind != yaml.MappingNode {
			return errors.New("invalid host entry")
		}

		var host Host
		if err := valNode.Decode(&host); err != nil {
			return err
		}

		if host.Host == "" {
			return errors.New("host entry missing host field for " + key)
		}

		hn.Hosts[key] = host
	}

	return nil
}

func (h *Host) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return errors.New("invalid host entry")
	}

	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valNode := value.Content[i+1]

		switch keyNode.Value {
		case "host":
			if valNode.Kind == yaml.ScalarNode {
				hostname := valNode.Value
				if strings.ContainsRune(valNode.Value, '@') {
					parts := strings.SplitN(valNode.Value, "@", 2)
					h.User = &parts[0]
					hostname = parts[1]
				}

				if strings.ContainsRune(hostname, ':') {
					parts := strings.SplitN(hostname, ":", 2)
					hostname = parts[0]
					if len(parts[1]) > 0 {
						port, err := strconv.Atoi(parts[1])
						if err != nil {
							return err
						}
						h.Port = &port
					}
				}

				h.Host = hostname
			}
		case "port":
			if valNode.Kind == yaml.ScalarNode {
				port, err := strconv.Atoi(valNode.Value)
				if err != nil {
					return err
				}
				h.Port = &port
			}
		case "user":
			if valNode.Kind == yaml.ScalarNode {
				h.User = &valNode.Value
			}
		case "identity":
			if valNode.Kind == yaml.ScalarNode {
				h.Identity = &valNode.Value
			}
		case "password":
			if valNode.Kind == yaml.ScalarNode {
				h.Password = &valNode.Value
			}
		case "groups":
			if valNode.Kind == yaml.SequenceNode {
				groups := []string{}
				for _, groupNode := range valNode.Content {
					if groupNode.Kind == yaml.ScalarNode {
						groups = append(groups, groupNode.Value)
					}
				}
				h.Groups = groups
			}

			if valNode.Kind == yaml.ScalarNode {
				h.Groups = []string{valNode.Value}
			}
		case "meta":
			if valNode.Kind == yaml.MappingNode {
				meta := map[string]interface{}{}
				if err := valNode.Decode(&meta); err != nil {
					return err
				}
				h.Meta = meta
			}
		case "os":
			if valNode.Kind == yaml.MappingNode {
				var os OS
				if err := valNode.Decode(&os); err != nil {
					return err
				}
				h.OS = &os
			} else if valNode.Kind == yaml.ScalarNode {
				os := &OS{
					Platform: valNode.Value,
				}
				h.OS = os
			} else {
				return errors.New("invalid host/os section")
			}
		}
	}

	return nil
}

func (hosts *Hosts) UnmarshalYAML(node *yaml.Node) error {

	for _, item := range node.Content {
		if item.Kind == yaml.ScalarNode {
			next := &Host{}
			hostname := item.Value
			if strings.ContainsRune(item.Value, '@') {
				parts := strings.SplitN(item.Value, "@", 2)
				next.User = &parts[0]
				hostname = parts[1]
			}

			if strings.ContainsRune(hostname, ':') {
				parts := strings.SplitN(hostname, ":", 2)
				hostname = parts[0]
				if len(parts[1]) > 0 {
					port, err := strconv.Atoi(parts[1])
					if err != nil {
						return err
					}
					next.Port = &port
				}
			}

			next.Host = hostname
			(*hosts)[hostname] = *next
			continue

		}

		if item.Kind != yaml.MappingNode {
			continue
		}

		for i := 0; i < len(item.Content); i += 2 {
			keyNode := item.Content[i]
			valNode := item.Content[i+1]

			if keyNode.Kind != yaml.ScalarNode {
				return errors.New("host key must be a string")
			}

			key := keyNode.Value

			if valNode.Kind == yaml.ScalarNode {
				next := &Host{}
				hostname := valNode.Value
				if strings.ContainsRune(valNode.Value, '@') {
					parts := strings.SplitN(valNode.Value, "@", 2)
					next.User = &parts[0]
					hostname = parts[1]
				}

				if strings.ContainsRune(hostname, ':') {
					parts := strings.SplitN(hostname, ":", 2)
					hostname = parts[0]
					if len(parts[1]) > 0 {
						port, err := strconv.Atoi(parts[1])
						if err != nil {
							return err
						}
						next.Port = &port
					}
				}

				next.Host = hostname
				(*hosts)[key] = *next
				continue
			}

			if valNode.Kind != yaml.MappingNode {
				return errors.New("invalid host entry")
			}

			var host Host
			if err := item.Decode(&host); err != nil {
				return err
			}

			if host.Host == "" {
				return errors.New("host entry missing host field for " + key)
			}

			(*hosts)[key] = host
		}
	}

	return nil
}
