module github.com/hyprxlabs/xtask

go 1.24.5

require (
	github.com/hyprxlabs/go/cmdargs v0.1.1
	github.com/hyprxlabs/go/dotenv v0.1.0
	github.com/hyprxlabs/go/env v0.1.4
	github.com/hyprxlabs/go/exec v0.1.4
	github.com/melbahja/goph v1.4.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.7
	golang.org/x/crypto v0.41.0
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/hyprxlabs/go/cmdargs v0.1.1 => ./xvendor/cmdargs
	github.com/hyprxlabs/go/dotenv v0.1.0 => ./xvendor/dotenv
	github.com/hyprxlabs/go/env v0.1.4 => ./xvendor/env
	github.com/hyprxlabs/go/exec v0.1.4 => ./xvendor/exec
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/elliotchance/orderedmap/v3 v3.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.13.9 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
