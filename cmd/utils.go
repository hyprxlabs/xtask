package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func getFile(file string, dir string) (string, error) {

	if file == "" && dir == "" {
		wd, _ := os.Getwd()
		if wd != "" {
			file = filepath.Join(wd, "xtaskfile")
			if _, err := os.Stat(file); err == nil {
				return file, nil
			}
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		file = filepath.Join(homeDir, "xtaskfile")
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
		return "", os.ErrNotExist
	}

	if file != "" {

		file, err := resolvePath(file)
		if err != nil {
			return "", err
		}

		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	if dir != "" {
		localTaskFile := ""
		cwd, err := os.Getwd()
		if err == nil {
			// fast check
			file = filepath.Join(cwd, dir, "xtaskfile")
			if _, err := os.Stat(file); err == nil {
				return file, nil
			}

			file = filepath.Join(cwd, "xtaskfile")
			if _, err := os.Stat(file); err == nil {
				localTaskFile = file
			}
		}

		if localTaskFile == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			file = filepath.Join(homeDir, "xtaskfile")
			if _, err := os.Stat(file); err == nil {
				localTaskFile = file
			}
		}

		if localTaskFile != "" {
			config := map[string]interface{}{}
			data, err := os.ReadFile(localTaskFile)
			if err != nil {
				return "", fmt.Errorf("error reading xtaskfile: %v", err)
			}

			if err := yaml.Unmarshal(data, &config); err != nil {
				return "", fmt.Errorf("error parsing xtaskfile: %v", err)
			}

			if config["config"] == nil {
				return "", fmt.Errorf("no config section found in xtaskfile")
			}

			configSection, ok := config["config"].(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("config section is not a mapping in xtaskfile")
			}

			dirs := []string{}
			wd, _ := os.Getwd()
			if len(wd) > 0 {
				dirs = append(dirs, wd)
			}

			obj, ok := configSection["delegated-dirs"]
			if !ok {
				return "", fmt.Errorf("delegated-dirs section is not defined in xtaskfile")
			}

			// determine what type of object it is
			switch obj := obj.(type) {
			case string:
				dirs = append(dirs, obj)

			case []interface{}:
				for _, v := range obj {
					if str, ok := v.(string); ok {
						dirs = append(dirs, str)
					}
				}
			}
			if len(dirs) == 0 {
				return "", fmt.Errorf("no directories found in xtaskfile config")
			}

			for _, d := range dirs {
				if !filepath.IsAbs(d) {
					n, err := filepath.Abs(d)
					if err != nil {
						return "", fmt.Errorf("error resolving directory: %v", err)
					}
					d = n
				}
				file = filepath.Join(d, dir, "xtaskfile")
				if _, err := os.Stat(file); err == nil {
					return file, nil
				}
			}

		}
	}

	return "", os.ErrNotExist
}

func resolvePath(file string) (string, error) {
	if file == "" {
		return os.Getwd()
	}

	if len(file) > 2 {
		if file[0] == '~' && (file[1] == '/' || file[1] == '\\') {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return filepath.Join(homeDir, file[2:], "xtaskfile"), nil
		} else if file[0] == '.' {
			i := 1
			if file[i] == '.' {
				i++
			}

			i++
			if file[i] == '/' || file[i] == '\\' {
				fp, err := filepath.Abs(file[i:])
				if err != nil {
					return "", err
				}
				return filepath.Join(fp, "xtaskfile"), nil
			}
		}
	}

	uri, _ := url.Parse(file)
	if uri != nil && uri.Scheme == "file://" && uri.Path != "" {
		return filepath.Clean(uri.Path), nil
	}

	if !filepath.IsAbs(file) {
		return filepath.Abs(file)
	}

	return file, nil
}
