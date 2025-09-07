package tasks

import (
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/hyprxlabs/go/env"
	"github.com/hyprxlabs/xtask/errors"
	"gopkg.in/yaml.v3"
)

func runTpl(ctx TaskContext) *TaskResult {

	uses := ctx.Data.Uses
	if uses == "tmpl" {
		uses = "tmpl://"
	}

	res := NewTaskResult()
	res.Start()

	uri, err := url.Parse(uses)
	if err != nil {
		return res.Fail(errors.New("Invalid template URI: " + err.Error()))
	}

	useTmpl := true
	useEnv := true

	if uri.Query().Get("disable-go-tmpl") == "true" {
		useTmpl = false
	}
	if uri.Query().Get("disable-env-tmpl") == "true" {
		useEnv = false
	}

	if env.Get("XTASK_DISABLE_GO_TMPL") == "true" {
		useTmpl = false
	}

	if env.Get("XTASK_DISABLE_ENV_TMPL") == "true" {
		useEnv = false
	}

	if uri.Scheme != "tmpl" {
		return res.Fail(errors.New("Invalid template URI scheme: " + uri.Scheme))
	}

	filesInput, ok := ctx.Data.With["files"]
	if !ok {
		return res.Fail(errors.New("No files to process for template task"))
	}

	filesArr, ok := filesInput.([]interface{})
	if !ok {
		return res.Fail(errors.New("Invalid files format for template task"))
	}

	files := []string{}
	for _, f := range filesArr {
		if str, ok := f.(string); ok {
			files = append(files, str)
		}
	}

	if files == nil {
		return res.Fail(errors.New("No files to process for template task"))
	}

	valuesFile := uri.Path
	values := map[string]interface{}{}
	if len(valuesFile) > 0 && valuesFile != "/" {
		bytes, err := os.ReadFile(valuesFile)
		if err != nil {
			return res.Fail(errors.New("Failed to read values file: " + err.Error()))
		}

		if err := yaml.Unmarshal(bytes, &values); err != nil {
			return res.Fail(errors.New("Failed to unmarshal values file: " + err.Error()))
		}
	}

	data := struct {
		env    map[string]string
		values map[string]interface{}
	}{
		env:    ctx.Data.Env.ToMap(),
		values: values,
	}

	for _, file := range files {
		src := file
		dest := file
		if strings.ContainsRune(file, ':') {
			parts := strings.SplitN(file, ":", 2)
			src = parts[0]
			dest = parts[1]
		} else {
			ext := filepath.Ext(src)
			switch ext {
			case ".tmpl", ".tpl", ".gotmpl":
				dest = strings.TrimSuffix(src, ext)
			default:
				dest = src + ".out"
			}
		}

		bytes, err := os.ReadFile(src)
		if err != nil {
			return res.Fail(errors.New("Failed to read template file: " + err.Error()))
		}
		content := string(bytes)
		if useEnv {
			updatedContent, err := env.ExpandWithOptions(content, &env.ExpandOptions{
				CommandSubstitution: true,
				Get: func(key string) string {
					if val, ok := data.env[key]; ok {
						return val
					}
					return ""
				},
				Set: func(key, value string) error {
					data.env[key] = value
					return nil
				},
			})

			if err != nil {
				return res.Fail(errors.New("Failed to expand environment variables in template: " + err.Error()))
			}

			content = updatedContent

			if !useTmpl {
				err := os.WriteFile(dest, []byte(content), 0644)
				if err != nil {
					return res.Fail(errors.New("Failed to write output file: " + err.Error()))
				}
			}
		}

		if useTmpl {
			tmp, err := template.New(src).Funcs(sprig.FuncMap()).Parse(content)
			if err != nil {
				return res.Fail(errors.New("Failed to parse template file: " + err.Error()))
			}

			out, err := os.Create(dest)
			if err != nil {
				return res.Fail(errors.New("Failed to create output file: " + err.Error()))
			}

			defer out.Close()

			if err := tmp.Execute(out, data); err != nil {
				return res.Fail(errors.New("Failed to execute template: " + err.Error()))
			}
		}
	}

	return res.Ok()
}
