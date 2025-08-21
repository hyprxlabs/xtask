package yamlenv

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// YamlToEnv converts YAML structure to environment variable names
type YamlToEnv struct {
	envVars map[string]string
}

// NewYamlToEnv creates a new YamlToEnv instance
func NewYamlToEnv() *YamlToEnv {
	return &YamlToEnv{
		envVars: make(map[string]string),
	}
}

// ProcessYAML reads YAML content and extracts environment variables
func (y *YamlToEnv) ProcessYAML(yamlContent []byte) error {
	var data interface{}
	err := yaml.Unmarshal(yamlContent, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	y.processNode(data, []string{})
	return nil
}

// processNode recursively processes YAML nodes
func (y *YamlToEnv) processNode(node interface{}, path []string) {
	switch v := node.(type) {
	case map[string]interface{}:
		// Handle map/object
		for key, value := range v {
			newPath := append(path, key)
			y.processNode(value, newPath)
		}
	case []interface{}:
		// Handle array - we'll index them
		for i, value := range v {
			newPath := append(path, fmt.Sprintf("%d", i))
			y.processNode(value, newPath)
		}
	default:
		// Handle leaf values (string, int, bool, etc.)
		if len(path) > 0 {
			envName := y.pathToEnvName(path)
			envValue := y.valueToString(v)
			y.envVars[envName] = envValue
		}
	}
}

// pathToEnvName converts a path array to environment variable name
func (y *YamlToEnv) pathToEnvName(path []string) string {
	// Convert to uppercase and join with underscores
	var parts []string
	for _, part := range path {
		// Convert to uppercase and replace non-alphanumeric with underscore
		part = strings.ToUpper(part)
		part = strings.ReplaceAll(part, "-", "_")
		part = strings.ReplaceAll(part, ".", "_")
		parts = append(parts, part)
	}
	return strings.Join(parts, "_")
}

// valueToString converts various types to string representation
func (y *YamlToEnv) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetEnvVars returns the extracted environment variables
func (y *YamlToEnv) GetEnvVars() map[string]string {
	return y.envVars
}

// PrintEnvVars prints all extracted environment variables
func (y *YamlToEnv) PrintEnvVars() {
	for name, value := range y.envVars {
		fmt.Printf("%s=%s\n", name, value)
	}
}

// SetEnvVars sets the environment variables in the current process
func (y *YamlToEnv) SetEnvVars() error {
	for name, value := range y.envVars {
		err := os.Setenv(name, value)
		if err != nil {
			return fmt.Errorf("failed to set env var %s: %w", name, err)
		}
	}
	return nil
}
