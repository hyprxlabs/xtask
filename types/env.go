package types

import (
	"iter"
	"maps"
	"os"
	"runtime"
	"strings"

	om "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
)

type Env struct {
	om      map[string]string
	keys    []string
	secrets []string
}

type EnvItem struct {
	Name   string
	Value  string
	Secret bool
}

func (e *Env) UnmarshalYAML(value *yaml.Node) error {
	if e == nil {
		e = NewEnv()
	}

	if value.Kind == yaml.SequenceNode {
		for _, item := range value.Content {
			if item.Kind == yaml.ScalarNode {
				parts := strings.SplitN(item.Value, "=", 2)
				if len(parts) == 2 {
					e.Set(parts[0], parts[1])
				} else {
					e.Set(parts[0], "")
				}
				continue
			}

			if item.Kind != yaml.MappingNode {
				continue
			}

			var envItem EnvItem
			if err := item.Decode(&envItem); err != nil {
				return err
			}
			e.Set(envItem.Name, envItem.Value)
			e.secrets = append(e.secrets, envItem.Name)
		}
	}

	if value.Kind == yaml.MappingNode {
		envMap := om.OrderedMap[string, string]{}
		if err := value.Decode(&envMap); err != nil {
			return err
		}
		for el := envMap.Oldest(); el != nil; el = el.Next() {
			e.Set(el.Key, el.Value)
		}
	}

	return nil
}

func NewEnv() *Env {
	return &Env{
		om:      map[string]string{},
		keys:    []string{},
		secrets: []string{},
	}
}

func NewEnvFromMap(omap om.OrderedMap[string, string]) *Env {
	keys := []string{}
	om := map[string]string{}
	for el := omap.Oldest(); el != nil; el = el.Next() {
		om[el.Key] = el.Value
		keys = append(keys, el.Key)
	}

	return &Env{
		om:      om,
		keys:    keys,
		secrets: []string{},
	}
}

func (e *Env) IsSecret(key string) bool {
	e.init()

	for _, k := range e.secrets {
		if k == key {
			return true
		}
	}
	return false
}

func (e *Env) Secrets() []string {
	e.init()
	if e.secrets == nil {
		e.secrets = []string{}
	}

	return e.secrets
}

func (e *Env) Set(key, value string) {
	e.init()

	if _, ok := e.om[key]; !ok {
		e.keys = append(e.keys, key)
	}

	e.om[key] = value
}

func (e *Env) Get(key string) (string, bool) {
	e.init()
	val, ok := e.om[key]
	return val, ok
}

func (e *Env) Has(key string) bool {
	e.init()
	_, ok := e.om[key]
	return ok
}

func (e *Env) PrependPath(path string) error {
	e.init()
	paths := e.SplitPath()

	if len(paths) > 0 {
		if runtime.GOOS == "windows" {
			if strings.EqualFold(paths[0], path) {
				return nil
			}
		} else {
			if paths[0] == path {
				return nil
			}
		}
	}

	paths = append([]string{path}, paths...)
	e.SetPath(strings.Join(paths, string(os.PathListSeparator)))
	return nil
}

func (e *Env) AppendPath(path string) error {
	e.init()
	paths := e.SplitPath()

	if len(paths) > 0 {
		if runtime.GOOS == "windows" {
			for _, p := range paths {
				if strings.EqualFold(p, path) {
					return nil
				}
			}
		} else {
			for _, p := range paths {
				if p == path {
					return nil
				}
			}
		}
	}

	paths = append(paths, path)
	e.SetPath(strings.Join(paths, string(os.PathListSeparator)))
	return nil
}

func (e *Env) HasPath(path string) bool {
	e.init()
	paths := e.SplitPath()
	if runtime.GOOS == "windows" {
		for _, p := range paths {
			if strings.EqualFold(p, path) {
				return true
			}
		}
		return false
	}

	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func (e *Env) SplitPath() []string {
	e.init()
	if e.GetPath() == "" {
		return []string{}
	}
	return strings.Split(e.GetPath(), string(os.PathListSeparator))
}

func (e *Env) GetPath() string {
	e.init()
	if runtime.GOOS == "windows" {
		if val, ok := e.om["Path"]; ok {
			return val
		}

		return ""
	}

	if val, ok := e.om["PATH"]; ok {
		return val
	}

	return ""
}

func (e *Env) SetPath(value string) error {
	e.init()
	if runtime.GOOS == "windows" {
		e.om["Path"] = value
		return nil
	}

	e.om["PATH"] = value
	return nil
}

func (e *Env) GetString(key string) string {
	e.init()
	if val, ok := e.om[key]; ok {
		return val
	}
	return ""
}

func (e *Env) Delete(key string) {
	e.init()
	delete(e.om, key)
	for i, k := range e.keys {
		if k == key {
			e.keys = append(e.keys[:i], e.keys[i+1:]...)
			break
		}
	}
}

func (e *Env) Clone() *Env {
	e.init()
	clone := NewEnv()

	for k, v := range e.om {
		clone.om[k] = v
	}
	clone.keys = append(clone.keys, e.keys...)
	clone.secrets = append(clone.secrets, e.secrets...)
	return clone
}

func (e *Env) ToOrderedMap() om.OrderedMap[string, string] {
	e.init()
	omap := om.New[string, string]()
	for _, k := range e.keys {
		omap.Set(k, e.om[k])
	}
	return *omap
}

func (e *Env) ToMap() map[string]string {
	e.init()
	m := make(map[string]string, len(e.om))
	maps.Copy(m, e.om)
	return m
}

func (e *Env) Keys() []string {
	e.init()
	keys := make([]string, 0, len(e.om))
	for k := range e.om {
		keys = append(keys, k)
	}
	return keys
}

func (e *Env) Values() []string {
	e.init()
	values := make([]string, 0, len(e.om))
	for _, k := range e.keys {
		values = append(values, e.om[k])
	}
	return values
}

func (e *Env) Len() int {
	e.init()
	return len(e.om)
}

// return iter.Seq
func (e *Env) Iter() iter.Seq2[string, string] {
	e.init()
	return func(yield func(string, string) bool) {
		for _, k := range e.keys {
			if !yield(k, e.om[k]) {
				break
			}
		}
	}
}

func (e *Env) init() {
	if e == nil {
		e = NewEnv()
	}

	if e.om == nil {
		e.om = map[string]string{}
	}

	if e.keys == nil {
		e.keys = []string{}
	}

	if e.secrets == nil {
		e.secrets = []string{}
	}
}
