package types

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Import struct {
	Uri       string `yaml:"uri"`
	Optional  bool   `yaml:"optional,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

type Imports []Import

func (i *Import) UnmarshalYAML(value *yaml.Node) error {

	if value.Kind == yaml.ScalarNode {
		i.Uri = value.Value
		return nil
	}

	if value.Kind != yaml.MappingNode {
		return errors.New("invalid import entry")
	}

	for j := 0; j < len(value.Content); j += 2 {
		keyNode := value.Content[j]
		valNode := value.Content[j+1]

		switch keyNode.Value {
		case "uri":
			if valNode.Kind == yaml.ScalarNode {
				i.Uri = valNode.Value
			}
		case "optional":
			if valNode.Kind == yaml.ScalarNode {
				i.Optional = valNode.Value == "true"
			}
		case "namespace":
			if valNode.Kind == yaml.ScalarNode {
				i.Namespace = valNode.Value
			}
		}
	}

	return nil
}
