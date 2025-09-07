package workflows

import (
	"errors"
	"os"

	"github.com/hyprxlabs/go/exec"
)

func (ws *Workflow) Exec(args []string) error {
	if ws == nil {
		return errors.New("workflow is nil")
	}

	dir, ok := ws.Env.Get("XTASK_DIR")
	if !ok || len(dir) == 0 {
		return errors.New("XTASK_DIR is not set")
	}

	cmd := exec.CommandContext(ws.Context, args[0]).WithArgs(args[1:]...)
	cmd.Dir = dir

	cmd.WithEnvMap(ws.Env.ToMap())
	o, err := cmd.Run()
	if err != nil || o.Code != 0 {
		if o.Code != 0 {
			os.Exit(o.Code)
		} else {
			os.Exit(1)
		}
	}

	os.Exit(o.Code)
	return nil
}
