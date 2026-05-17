package bdd

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// --- block.feature support ------------------------------------------------

func initializeBlockSteps(ctx *godog.ScenarioContext, w *world) {
	ctx.Step(`^"([^"]*)" has a note containing "([^"]*)"$`, w.taskHasNoteContaining)
	ctx.Step(`^the output reports that a reason is required$`, w.outputReportsReasonRequired)
	ctx.Step(`^the command reports the task is not blocked$`, w.outputReportsTaskNotBlocked)
}

func (w *world) taskHasNoteContaining(id, text string) error {
	t, err := loadFixtureTask(id)
	if err != nil {
		return err
	}
	idx := strings.Index(t.Body, "## Notes")
	if idx < 0 {
		return fmt.Errorf("task %s has no ## Notes section; body:\n%s", id, t.Body)
	}
	if !strings.Contains(t.Body[idx:], text) {
		return fmt.Errorf("note for %s does not contain %q; body:\n%s", id, text, t.Body)
	}
	return nil
}

func (w *world) outputReportsReasonRequired() error {
	combined := w.stdout.String() + w.stderr.String()
	if !strings.Contains(combined, "reason is required") {
		return fmt.Errorf("expected output to report reason required; got: %s", combined)
	}
	return nil
}

func (w *world) outputReportsTaskNotBlocked() error {
	combined := w.stdout.String() + w.stderr.String()
	if !strings.Contains(combined, "not blocked") {
		return fmt.Errorf("expected output to report task not blocked; got: %s", combined)
	}
	return nil
}
