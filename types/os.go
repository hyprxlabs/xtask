package types

type OS struct {
	Platform     string `yaml:"platform"`
	Arch         string `yaml:"arch"`
	Variant      string `yaml:"variant,omitempty"`
	Family       string `yaml:"family,omitempty"`
	Codename     string `yaml:"codename,omitempty"`
	Version      string `yaml:"version,omitempty"`
	BuildVersion string `yaml:"build_version,omitempty"`
}
