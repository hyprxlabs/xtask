package workflows

import "github.com/hyprxlabs/xtask/types"

type CyclicalReferenceError struct {
	Cycles []types.Task
}

func (e *CyclicalReferenceError) Error() string {
	msg := "Cyclical references found in tasks:\n"
	for _, cycle := range e.Cycles {
		msg += " - " + cycle.Id + "\n"
	}
	return msg
}
