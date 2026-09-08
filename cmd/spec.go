package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aholbreich/tl/internal/spec"
	"github.com/aholbreich/tl/internal/task"
)

// specStatusFlagUsage is shared so list and ready describe the flag
// identically.
const specStatusFlagUsage = "Resolve referenced .feature specs (reads those files)"

// resolveSpecsFor reads the specifications referenced by tasks. It returns nil
// when enabled is false, which is what makes the JSON "spec" key null unless
// the flag was passed — a consumer's schema never depends on the flag, only
// the value does.
//
// Reading happens here and nowhere else: the default path of every read
// command must not open a file outside .tl/ (decision 0002).
func resolveSpecsFor(ledger string, tasks []*task.Task, enabled bool) map[string][]spec.Spec {
	if !enabled {
		return nil
	}
	repoRoot := filepath.Dir(ledger)
	out := make(map[string][]spec.Spec, len(tasks))
	for _, t := range tasks {
		// Always non-nil, so a task with no spec reference emits [] rather
		// than null once the flag is on.
		out[t.ID] = spec.Resolve(repoRoot, t.References)
	}
	return out
}

// specCell renders one task's specs for a table column.
func specCell(specs []spec.Spec) string {
	if len(specs) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(specs))
	for _, s := range specs {
		parts = append(parts, s.Path+" ("+specDetail(s)+")")
	}
	return strings.Join(parts, ", ")
}

func specDetail(s spec.Spec) string {
	switch s.State {
	case spec.StatePresent:
		return fmt.Sprintf("%d", s.Scenarios)
	default:
		return string(s.State)
	}
}
