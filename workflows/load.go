package workflows

import (
	"errors"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hyprxlabs/go/dotenv"
	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/paths"
	"github.com/hyprxlabs/xtask/types"
	"github.com/hyprxlabs/xtask/versions"
	"gopkg.in/yaml.v3"
)

func (wf *Workflow) Load(taskfile types.XTaskfile) error {

	err := wf.LoadEnv(taskfile)
	if err != nil {
		return err
	}
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldDir)
	rootDir := filepath.Dir(taskfile.Path)
	os.Chdir(rootDir)

	hosts := map[string]types.Host{}
	envMap := wf.Env
	if len(taskfile.HostsNode.Imports) > 0 {
		imports := taskfile.HostsNode.Imports
		keys := envMap.Keys()
		for _, imp := range imports {
			opts := &env.ExpandOptions{
				Get: func(key string) string {
					s, ok := envMap.Get(key)
					if ok {
						return s
					}

					return ""
				},
				Set: func(key, value string) error {
					envMap.Set(key, value)
					return nil
				},
				Keys:                keys,
				ExpandUnixArgs:      true,
				ExpandWindowsVars:   false,
				CommandSubstitution: taskfile.Config.Substitution,
			}
			next, err := env.ExpandWithOptions(imp, opts)
			if err != nil {
				return errors.New("failed to expand hosts import path: " + imp + " error: " + err.Error())
			}

			next = strings.TrimSpace(next)
			optional := false
			if strings.HasSuffix(next, "?") {
				optional = true
				next = strings.TrimSuffix(next, "?")
			}

			if !(filepath.IsAbs(next)) {
				p, err := filepath.Abs(next)
				if err != nil {
					return errors.New("failed to get absolute path of hosts import: " + next + " error: " + err.Error())
				}
				next = p
			}

			if !isFile(next) {
				if optional {
					continue
				} else {
					return errors.New("required hosts import file does not exist: " + next)
				}
			}

			data, err := os.ReadFile(next)
			if err != nil {
				return errors.New("failed to read hosts import file: " + next + " error: " + err.Error())
			}

			var hostfile types.XHostFile
			err = hostfile.Decode(data)
			if err != nil {
				return errors.New("failed to parse hosts import file: " + next + " error: " + err.Error())
			}

			maps.Copy(hosts, hostfile.Hosts)
		}
	}

	wf.Hosts = hosts
	keys := envMap.Keys()
	for k, v := range *taskfile.Tasks {
		wf.Tasks[k] = v

		run := ""
		if v.Run != nil && len(*v.Run) > 0 {
			run = *v.Run
		}

		run = strings.TrimSpace(run)

		if len(run) > 0 {

			if !strings.ContainsAny(run, "\n\r") && (strings.HasSuffix(run, ".task.yaml") || strings.HasSuffix(run, ".task.yml")) {
				if strings.Contains(run, "://") {
					uri, err := url.Parse(run)
					if err != nil {
						return errors.New("failed to parse task import URI: " + run + " error: " + err.Error())
					}

					if uri.Scheme != "file" {
						continue
					}

					r := uri.Path
					run = strings.TrimSpace(r)
				}

				opts := &env.ExpandOptions{
					Get: func(key string) string {
						s, ok := envMap.Get(key)
						if ok {
							return s
						}

						return ""
					},
					Set: func(key, value string) error {
						envMap.Set(key, value)
						return nil
					},
					Keys:                keys,
					ExpandUnixArgs:      true,
					ExpandWindowsVars:   false,
					CommandSubstitution: taskfile.Config.Substitution,
				}

				next, err := env.ExpandWithOptions(run, opts)
				run = next
				if err != nil {
					return errors.New("failed to expand task import path: " + run + " error: " + err.Error())
				}

				run = strings.TrimSpace(run)
				if !(filepath.IsAbs(run)) {
					p, err := filepath.Abs(filepath.Join(rootDir, run))
					if err != nil {
						return errors.New("failed to get absolute path of task import: " + run + " error: " + err.Error())
					}
					run = p
				}

				if !isFile(run) {
					return errors.New("task import file does not exist: " + run)
				}

				data, err := os.ReadFile(run)
				if err != nil {
					return errors.New("failed to read task import file: " + run + " error: " + err.Error())
				}

				var task types.Task
				err = yaml.Unmarshal(data, &task)
				if err != nil {
					return errors.New("failed to parse task import file: " + run + " error: " + err.Error())
				}

				wf.Tasks[k] = task
			}
		}
	}

	// continue loadding other parts like Hosts, Tasks, etc.

	return nil
}

func (wf *Workflow) LoadEnv(taskfile types.XTaskfile) error {

	if len(taskfile.Path) == 0 {
		return errors.New("taskfile path is empty")
	}

	if !(filepath.IsAbs(taskfile.Path)) {
		p, err := filepath.Abs(taskfile.Path)
		if err != nil {
			return err
		}
		taskfile.Path = p
	}

	rootDir := filepath.Dir(taskfile.Path)
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	os.Chdir(rootDir)
	defer os.Chdir(currentDir)

	if wf == nil {
		wf = NewWorkflow()
	}

	envMap := types.NewEnv()
	for _, n := range os.Environ() {
		parts := strings.SplitN(n, "=", 2)
		if len(parts) == 2 {
			envMap.Set(parts[0], parts[1])
		} else {
			envMap.Set(parts[0], "")
		}
	}

	normalizeEnv(envMap)
	envMap.Set("XTASK_FILE", taskfile.Path)
	envMap.Set("XTASK_DIR", rootDir)
	envMap.Set("XTASK_CONTEXT", wf.ContextName)
	envMap.Set("XTASK_SHELL", taskfile.Config.Shell)
	envMap.Set("XTASK_ETC_DIR", wf.Config.Dirs.Etc)
	configHome := envMap.GetString("XTASK_CONFIG_HOME")
	if configHome == "" {
		configHome, _ = paths.UserConfigDir()
	}
	envMap.Set("XTASK_CONFIG_HOME", configHome)
	dataHome := envMap.GetString("XTASK_DATA_HOME")
	if dataHome == "" {
		dataHome, _ = paths.UserDataDir()
	}
	envMap.Set("XTASK_DATA_HOME", dataHome)
	cacheHome := envMap.GetString("XTASK_CACHE_HOME")
	if cacheHome == "" {
		cacheHome, _ = paths.UserCacheDir()
	}
	envMap.Set("XTASK_CACHE_HOME", cacheHome)
	stateHome := envMap.GetString("XTASK_STATE_HOME")
	if stateHome == "" {
		stateHome, _ = paths.UserStateDir()
	}
	envMap.Set("XTASK_STATE_HOME", stateHome)
	envMap.Set("XTASK_APPS_DIRS", strings.Join(wf.Config.Dirs.Apps, string(os.PathListSeparator)))
	envMap.Set("XTASK_VERSION", versions.Version) // TODO: set actual version

	if _, ok := envMap.Get("XTASK_ENV"); !ok {
		f, err := os.CreateTemp("", "xtask-env-")
		if err != nil {
			return err
		}
		f.Write([]byte{})
		f.Close()
		envMap.Set("XTASK_ENV", f.Name())
		wf.cleanupEnv = true
	}

	if _, ok := envMap.Get("XTASK_PATH"); !ok {
		f, err := os.CreateTemp("", "xtask-path-")
		if err != nil {
			return err
		}
		f.Write([]byte{})
		f.Close()

		envMap.Set("XTASK_PATH", f.Name())
		wf.cleanupPath = true
	}

	if len(taskfile.Config.PrependPaths) > 0 {
		for _, p := range taskfile.Config.PrependPaths {
			if p.OS != "" {
				if !strings.EqualFold(p.OS, runtime.GOOS) {
					continue
				}
			}

			opts := &env.ExpandOptions{
				Get: func(key string) string {
					s, ok := envMap.Get(key)
					if ok {
						return s
					}

					return ""
				},
				Set: func(key, value string) error {
					envMap.Set(key, value)
					return nil
				},
				Keys:                envMap.Keys(),
				ExpandUnixArgs:      true,
				ExpandWindowsVars:   false,
				CommandSubstitution: taskfile.Config.Substitution,
			}
			path, err := env.ExpandWithOptions(p.Path, opts)
			if err != nil {
				return errors.New("failed to expand prepend-path: " + p.Path + " error: " + err.Error())
			}
			path = strings.TrimSpace(path)
			if !(filepath.IsAbs(path)) {
				abs, err := filepath.Abs(filepath.Join(rootDir, path))
				if err != nil {
					return errors.New("failed to get absolute path of prepend-path: " + path + " error: " + err.Error())
				}
				path = abs
			}

			env.PrependPath(path)
		}
	}

	dotenvFiles := []string{}
	if wf.Config.Dirs.Etc == "" {
		wf.Config.Dirs.Etc = "./.xtask/etc"
	}

	if len(wf.Config.Dirs.Apps) == 0 {
		wf.Config.Dirs.Apps = []string{"./.xtask/apps/*"}
		data := envMap.GetString("XTASK_DATA_HOME")
		if len(data) > 0 {
			wf.Config.Dirs.Apps = append(wf.Config.Dirs.Apps, data+"/apps/*")
		}
	}

	if isDir(configHome) {
		dotenvFiles = append(dotenvFiles, filepath.Join(configHome, ".env.shared?"))
		if wf.ContextName == "default" || wf.ContextName == "" {
			dotenvFiles = append(dotenvFiles, filepath.Join(configHome, ".env?"))
			dotenvFiles = append(dotenvFiles, filepath.Join(configHome, ".env.default?"))
		} else {
			dotenvFiles = append(dotenvFiles, filepath.Join(configHome, ".env."+wf.ContextName+"?"))
		}
	}

	if isDir(wf.Config.Dirs.Etc) {
		dotenvFiles = append(dotenvFiles, filepath.Join(wf.Config.Dirs.Etc, ".env.shared?"))
		if wf.ContextName == "default" || wf.ContextName == "" {
			dotenvFiles = append(dotenvFiles, filepath.Join(wf.Config.Dirs.Etc, ".env?"))
			dotenvFiles = append(dotenvFiles, filepath.Join(wf.Config.Dirs.Etc, ".env.default?"))
		} else {
			dotenvFiles = append(dotenvFiles, filepath.Join(wf.Config.Dirs.Etc, ".env."+wf.ContextName+"?"))
		}
	}

	if wf.ContextName == "default" || wf.ContextName == "" {
		dotenvFiles = append(dotenvFiles, filepath.Join(rootDir, ".env?"))
		dotenvFiles = append(dotenvFiles, filepath.Join(rootDir, ".env.default?"))
	} else {
		dotenvFiles = append(dotenvFiles, filepath.Join(rootDir, ".env."+wf.ContextName+"?"))
	}

	if len(taskfile.Dotenv) > 0 {
		for _, f := range taskfile.Dotenv {
			skip := false
			for _, existing := range dotenvFiles {
				if existing == f {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
			dotenvFiles = append(dotenvFiles, f)
		}
	}

	if taskfile.Config.Env.Len() > 0 {
		keys := envMap.Keys()
		opts := &env.ExpandOptions{
			Get: func(key string) string {
				s, ok := envMap.Get(key)
				if ok {
					return s
				}

				return ""
			},
			Set: func(key, value string) error {
				envMap.Set(key, value)
				return nil
			},
			Keys:                keys,
			ExpandUnixArgs:      true,
			ExpandWindowsVars:   false,
			CommandSubstitution: taskfile.Config.Substitution,
		}

		for k, v := range taskfile.Config.Env.Iter() {
			expandedValue, err := env.ExpandWithOptions(v, opts)
			if err != nil {
				return err
			}
			envMap.Set(k, expandedValue)

			hasKey := false
			for _, key := range opts.Keys {
				if key == k {
					hasKey = true
					break
				}
			}
			if !hasKey {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	if len(dotenvFiles) > 0 {
		opts := &env.ExpandOptions{
			Get: func(key string) string {
				s, ok := envMap.Get(key)
				if ok {
					return s
				}

				return ""
			},
			Set: func(key, value string) error {
				envMap.Set(key, value)
				return nil
			},
			Keys:                envMap.Keys(),
			ExpandUnixArgs:      true,
			ExpandWindowsVars:   false,
			CommandSubstitution: taskfile.Config.Substitution,
		}

		globalDoc := dotenv.NewDocument()

		for _, f := range dotenvFiles {
			next, err := env.ExpandWithOptions(f, opts)
			if err != nil {
				return errors.New("failed to expand dotenv file path: " + f + " error: " + err.Error())
			}

			optional := false
			if strings.HasSuffix(next, "?") {
				optional = true
				next = strings.TrimSuffix(next, "?")
			}

			if !(filepath.IsAbs(next)) {
				abs, err := filepath.Abs(next)
				if err != nil {
					return errors.New("failed to get absolute path of dotenv file: " + next + " error: " + err.Error())
				}
				next = abs
			}

			if !isFile(next) {
				if optional {
					continue
				} else {
					return errors.New("required dotenv file does not exist: " + next)
				}
			}

			data, err := os.ReadFile(next)
			if err != nil {
				return errors.New("failed to read dotenv file: " + next + " error: " + err.Error())
			}

			doc, err := dotenv.Parse(string(data))
			if err != nil {
				return errors.New("failed to parse dotenv file: " + next + " error: " + err.Error())
			}

			globalDoc.Merge(doc)
		}

		for _, node := range globalDoc.ToArray() {
			if node.Type != dotenv.VARIABLE_TOKEN {
				continue
			}

			keyPtr := node.Key
			if keyPtr == nil {
				continue
			}

			key := *keyPtr
			value := node.Value
			expandedValue, err := env.ExpandWithOptions(value, opts)
			if err != nil {
				return errors.New("failed to expand dotenv variable: " + key + " error: " + err.Error())
			}
			envMap.Set(key, expandedValue)

			hasKey := false
			for _, k := range opts.Keys {
				if k == key {
					hasKey = true
					break
				}
			}
			if !hasKey {
				opts.Keys = append(opts.Keys, key)
			}
		}
	}

	if taskfile.Env.Len() > 0 {
		opts := &env.ExpandOptions{
			Get: func(key string) string {
				s, ok := envMap.Get(key)
				if ok {
					return s
				}

				return ""
			},
			Set: func(key, value string) error {
				envMap.Set(key, value)
				return nil
			},
			Keys:                envMap.Keys(),
			ExpandUnixArgs:      true,
			ExpandWindowsVars:   false,
			CommandSubstitution: taskfile.Config.Substitution,
		}
		for k, v := range taskfile.Env.Iter() {
			expandedValue, err := env.ExpandWithOptions(v, opts)
			if err != nil {
				return err
			}
			envMap.Set(k, expandedValue)

			hasKey := false
			for _, key := range opts.Keys {
				if key == k {
					hasKey = true
					break
				}
			}
			if !hasKey {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	wf.Env = envMap

	return nil
}

func normalizeEnv(envMap *types.Env) error {

	configHome := os.Getenv("XDG_CONFIG_HOME")
	dataHome := os.Getenv("XDG_DATA_HOME")
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	stateHome := os.Getenv("XDG_STATE_HOME")
	binHome := os.Getenv("XDG_BIN_HOME")
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	envMap.Set("OS_PLATFORM", runtime.GOOS)
	envMap.Set("OS_ARCH", runtime.GOARCH)

	if runtime.GOOS == "windows" {
		envMap.Set("OSTYPE", "windows")
		user, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		host, err := os.Hostname()
		if err != nil {
			return err
		}
		shell, ok := os.LookupEnv("SHELL")
		if !ok {
			shell = "powershell.exe"
		}
		envMap.Set("HOME", user)
		envMap.Set("HOMEPATH", user)
		envMap.Set("USER", user)
		envMap.Set("HOSTNAME", host)
		envMap.Set("SHELL", shell)

		if len(configHome) == 0 {
			configHome = filepath.Join(user, "AppData", "Roaming")
			envMap.Set("XDG_CONFIG_HOME", configHome)
		}

		if len(dataHome) == 0 {
			dataHome = filepath.Join(user, "AppData", "Local")
			envMap.Set("XDG_DATA_HOME", dataHome)
		}

		if len(cacheHome) == 0 {
			cacheHome = filepath.Join(user, "AppData", "Local", "Cache")
			envMap.Set("XDG_CACHE_HOME", cacheHome)
		}

		if len(stateHome) == 0 {
			stateHome = filepath.Join(user, "AppData", "Local", "State")
			envMap.Set("XDG_STATE_HOME", stateHome)
		}

		if len(binHome) == 0 {
			binHome = filepath.Join(user, "AppData", "Local", "Programs", "bin")
			envMap.Set("XDG_BIN_HOME", binHome)
		}

		if len(runtimeDir) == 0 {
			runtimeDir = filepath.Join(user, "AppData", "Local", "Temp")
			envMap.Set("XDG_RUNTIME_DIR", runtimeDir)
		}
	} else {
		osType := os.Getenv("OSTYPE")
		if len(osType) == 0 {
			osType = runtime.GOOS
			envMap.Set("OSTYPE", osType)
		}

		user, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		if len(configHome) == 0 {
			configHome = filepath.Join(user, ".config")
			envMap.Set("XDG_CONFIG_HOME", configHome)
		}

		if len(dataHome) == 0 {
			dataHome = filepath.Join(user, ".local", "share")
			envMap.Set("XDG_DATA_HOME", dataHome)
		}

		if len(cacheHome) == 0 {
			cacheHome = filepath.Join(user, ".cache")
			envMap.Set("XDG_CACHE_HOME", cacheHome)
		}

		if len(stateHome) == 0 {
			stateHome = filepath.Join(user, ".local", "state")
			envMap.Set("XDG_STATE_HOME", stateHome)
		}

		if len(binHome) == 0 {
			binHome = filepath.Join(user, ".local", "bin")
			envMap.Set("XDG_BIN_HOME", binHome)
		}

		if len(runtimeDir) == 0 {
			id := os.Getuid()
			runtimeDir = filepath.Join("user", "run", fmt.Sprintf("%d", id))
			envMap.Set("XDG_RUNTIME_DIR", runtimeDir)
		}
	}

	return nil
}
