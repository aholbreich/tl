package bdd

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func initializeDashboardSteps(ctx *godog.ScenarioContext, w *world) {
	ctx.Step(`^the dashboard Markdown is exactly:$`, w.dashboardMarkdownIsExactly)
	ctx.Step(`^the dashboard contains the line "([^"]*)"$`, w.dashboardContainsLine)
}

func (w *world) dashboardMarkdownIsExactly(expected *godog.DocString) error {
	if w.cmdErr != nil {
		return w.cmdErr
	}
	if got := w.stdout.String(); got != expected.Content+"\n" {
		return fmt.Errorf("dashboard = %q, want %q", got, expected.Content+"\n")
	}
	return nil
}

func (w *world) dashboardContainsLine(expected string) error {
	if w.cmdErr != nil {
		return w.cmdErr
	}
	for _, line := range strings.Split(w.stdout.String(), "\n") {
		if line == expected {
			return nil
		}
	}
	return fmt.Errorf("dashboard missing line %q; got:\n%s", expected, w.stdout.String())
}
