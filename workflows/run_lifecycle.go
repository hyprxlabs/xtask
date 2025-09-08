package workflows

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hyprxlabs/xtask/types"
)

func (wf *Workflow) RunLifecycle(target string, app string, contextName string) error {

	find := map[string]string{}
	for key := range wf.Tasks {
		if key == "test" || strings.HasSuffix(key, "test:") {
			find[key] = key
		}
	}
	targets := []string{}
	if app == "" || app == "default" {
		testFound := false
		if _, ok := wf.Tasks[target+":default:"+contextName+":before"]; ok {
			targets = append(targets, target+":default:"+contextName+":before")
		} else if _, ok := wf.Tasks[target+":default:before"]; ok {
			targets = append(targets, target+":default:before")
		} else if _, ok := wf.Tasks[target+":before"]; ok {
			targets = append(targets, target+":before")
		}

		if _, ok := wf.Tasks[target+":default:"+contextName]; ok {
			targets = append(targets, target+":default:"+contextName)
			testFound = true
		} else if _, ok := wf.Tasks[target+":default"]; ok {
			targets = append(targets, target+":default")
			testFound = true
		} else if _, ok := wf.Tasks[target]; ok {
			targets = append(targets, target)
			testFound = true
		}

		if _, ok := wf.Tasks[target+":default:"+contextName+":after"]; ok {
			targets = append(targets, target+":default:"+contextName+":after")
		} else if _, ok := wf.Tasks[target+":default:after"]; ok {
			targets = append(targets, target+":default:after")
		} else if _, ok := wf.Tasks[target+":after"]; ok {
			targets = append(targets, target+":after")
		}

		if !testFound {
			return errors.New("no default test task found")
		}

	} else {
		testFound := false

		if _, ok := find[target+":"+app+":"+contextName+":before"]; ok {
			targets = append(targets, target+":"+app+":"+contextName+":before")
		} else if _, ok := wf.Tasks[target+":"+app+":before"]; ok {
			targets = append(targets, target+":"+app+":before")
		} else if _, ok := wf.Tasks[target+":before"]; ok {
			targets = append(targets, target+":before")
		}

		if _, ok := find[target+":"+app+":"+contextName]; ok {
			targets = append(targets, target+":"+app+":"+contextName)
			testFound = true
		} else if _, ok := wf.Tasks[target+":"+app]; ok {
			targets = append(targets, target+":"+app)
			testFound = true
		}

		if _, ok := wf.Tasks[target+":"+app+":"+contextName+":after"]; ok {
			targets = append(targets, target+":"+app+":"+contextName+":after")
		} else if _, ok := wf.Tasks[target+":"+app+":after"]; ok {
			targets = append(targets, target+":"+app+":after")
		} else if _, ok := wf.Tasks[target+":after"]; ok {
			targets = append(targets, target+":after")
		}

		if wf.parent != nil && !testFound {
			return errors.New("no " + target + " task found for app: " + app + " in workflow")
		}

		if !testFound {
			apps := wf.Config.Dirs.Apps
			slices.Reverse(apps)

			og, err := os.Getwd()
			if err != nil {
				return err
			}

			baseDir := wf.Env.GetString("XTASK_DIR")
			if baseDir != "" && baseDir != og {
				os.Chdir(baseDir)
				defer os.Chdir(og)
			}

			nextTaskfile := ""

			for _, dir := range apps {
				dir = strings.TrimSpace(strings.TrimRight(dir, "*"))

				if !filepath.IsAbs(dir) {
					resolved, err := filepath.Abs(dir)
					if err != nil {
						return err
					}
					dir = resolved
				}

				basename := filepath.Base(dir)
				if strings.EqualFold(basename, app) {
					try := filepath.Join(dir, contextName, "xtaskfile")
					if isFile(try) {
						nextTaskfile = try
						break
					}

					try = filepath.Join(dir, "xtaskfile")
					if isFile(try) {
						nextTaskfile = try
						break
					}
				}

				try := filepath.Join(dir, app, "xtaskfile")
				if isFile(try) {
					nextTaskfile = try
					break
				}
			}

			if nextTaskfile != "" {
				tf := types.NewXTaskfile()
				err := tf.DecodeYAMLFile(nextTaskfile)
				if err != nil {
					return errors.New("Failed to read xtaskfile: " + nextTaskfile + " " + err.Error())
				}

				for _, k := range wf.Env.Keys() {
					v0 := wf.Env.GetString(k)
					if v, ok := os.LookupEnv(k); !ok || v0 != v {
						os.Setenv(k, v0)
					}
				}

				wf2 := NewWorkflow()
				err = wf2.Load(*tf)

				if err != nil {
					return errors.New("Failed to load xtaskfile: " + nextTaskfile + " " + err.Error())
				}

				wf2.parent = wf
				// app must be empty
				return wf2.RunLifecycle(target, "", contextName)
			}
		}
	}

	if len(targets) == 0 {
		return errors.New("no test tasks found")
	}

	wf.ContextName = contextName
	return wf.Run(targets, []string{})
}
