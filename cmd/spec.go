package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aholbreich/tl/internal/spec"
	"github.com/aholbreich/tl/internal/task"
)

// resolveSpecsFor reads the specifications referenced by tasks.
//
// It is unconditional. A task with no .feature reference resolves to an empty
// slice and costs nothing, so the work is already proportional to how much
// Gherkin a project actually writes — measured at 3ms across 500 spec files,
// which is why the flag that used to gate this was removed (decision 0002,
// amended).
func resolveSpecsFor(ledger string, tasks []*task.Task) map[string][]spec.Spec {
	repoRoot := filepath.Dir(ledger)
	out := make(map[string][]spec.Spec, len(tasks))
	for _, t := range tasks {
		// Always non-nil, so a task with no spec reference emits [] rather
		// than null.
		out[t.ID] = spec.Resolve(repoRoot, t.References)
	}
	return out
}

// anySpecs reports whether any task carries a spec reference. Human-facing
// tables use it to decide whether the column is worth its width: a project
// that writes no Gherkin sees exactly the output it saw before this existed.
func anySpecs(specs map[string][]spec.Spec) bool {
	for _, list := range specs {
		if len(list) > 0 {
			return true
		}
	}
	return false
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
