package types

type Need struct {
	Name     string `yaml:"name"`
	Parallel bool   `yaml:"async,omitempty"`
}
