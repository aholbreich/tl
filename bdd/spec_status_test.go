package bdd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// --- spec-status.feature support ------------------------------------------

func initializeSpecStatusSteps(ctx *godog.ScenarioContext, w *world) {
	// Fixtures
	ctx.Step(`^the repository has a feature file "([^"]*)" with (\d+) scenarios?$`, w.writeFeatureFile)
	ctx.Step(`^the repository has a feature file "([^"]*)" with (\d+) scenario and (\d+) scenario outline of (\d+) examples$`, w.writeFeatureFileWithOutline)
	ctx.Step(`^the repository has an unreadable file "([^"]*)"$`, w.writeUnreadableFile)

	// Table assertions
	ctx.Step(`^the output lists "([^"]*)" with spec "([^"]*)" and (\d+) scenarios$`, w.outputListsSpecWithScenarios)
	ctx.Step(`^the output lists "([^"]*)" with spec "([^"]*)" marked (missing|unknown)$`, w.outputListsSpecMarked)
	ctx.Step(`^the output lists "([^"]*)" with no spec$`, w.outputListsNoSpec)
	ctx.Step(`^the output lists "([^"]*)" with (\d+) specs$`, w.outputListsSpecCount)

	// JSON assertions
	ctx.Step(`^the JSON task "([^"]*)" has a null "spec" field$`, w.jsonTaskSpecNull)
	ctx.Step(`^the JSON task "([^"]*)" has an empty "spec" array$`, w.jsonTaskSpecEmpty)
	ctx.Step(`^the JSON task "([^"]*)" has a spec entry for "([^"]*)" that exists with (\d+) scenarios$`, w.jsonTaskSpecEntry)
}

// --- fixtures -------------------------------------------------------------

// writeFeatureFile writes a syntactically real feature file, so the counter is
// exercised against Gherkin rather than against a line of bare keywords.
func (w *world) writeFeatureFile(path string, count int) error {
	var b strings.Builder
	b.WriteString("@implemented\nFeature: Generated fixture\n")
	for i := 0; i < count; i++ {
		fmt.Fprintf(&b, "\n  Scenario: Generated scenario %d\n    Given a precondition\n    Then an outcome\n", i+1)
	}
	return writeRepoFile(path, b.String())
}

func (w *world) writeFeatureFileWithOutline(path string, scenarios, outlines, examples int) error {
	var b strings.Builder
	b.WriteString("@implemented\nFeature: Generated fixture\n")
	for i := 0; i < scenarios; i++ {
		fmt.Fprintf(&b, "\n  Scenario: Generated scenario %d\n    Given a precondition\n    Then an outcome\n", i+1)
	}
	for i := 0; i < outlines; i++ {
		fmt.Fprintf(&b, "\n  Scenario Outline: Generated outline %d\n    Given a precondition of <value>\n    Then an outcome\n\n    Examples:\n      | value |\n", i+1)
		for j := 0; j < examples; j++ {
			fmt.Fprintf(&b, "      | %d |\n", j+1)
		}
	}
	return writeRepoFile(path, b.String())
}

func (w *world) writeUnreadableFile(path string) error {
	if err := writeRepoFile(path, "@implemented\nFeature: Unreadable\n"); err != nil {
		return err
	}
	return os.Chmod(path, 0o000)
}

func writeRepoFile(path, content string) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// --- table assertions -----------------------------------------------------

func (w *world) outputListsSpecWithScenarios(id, path string, count int) error {
	return w.specCellContains(id, fmt.Sprintf("%s (%d)", path, count))
}

func (w *world) outputListsSpecMarked(id, path, state string) error {
	return w.specCellContains(id, fmt.Sprintf("%s (%s)", path, state))
}

func (w *world) outputListsNoSpec(id string) error {
	line, err := w.lineForTask(id)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimRight(line, " "), "-") {
		return fmt.Errorf("expected no spec for %s, got line: %s", id, line)
	}
	return nil
}

func (w *world) outputListsSpecCount(id string, count int) error {
	line, err := w.lineForTask(id)
	if err != nil {
		return err
	}
	got := strings.Count(line, ".feature")
	if got != count {
		return fmt.Errorf("expected %d specs for %s, found %d in line: %s", count, id, got, line)
	}
	return nil
}

func (w *world) specCellContains(id, want string) error {
	line, err := w.lineForTask(id)
	if err != nil {
		return err
	}
	if !strings.Contains(line, want) {
		return fmt.Errorf("expected %q in the row for %s; got: %s", want, id, line)
	}
	return nil
}

func (w *world) lineForTask(id string) (string, error) {
	for _, line := range strings.Split(w.stdout.String(), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), id) {
			return line, nil
		}
	}
	return "", fmt.Errorf("output has no row for %s; got:\n%s", id, w.stdout.String())
}

// --- JSON assertions ------------------------------------------------------

func (w *world) jsonTaskSpecNull(id string) error {
	raw, err := w.jsonSpecField(id)
	if err != nil {
		return err
	}
	if string(raw) != "null" {
		return fmt.Errorf("JSON task %s spec = %s, expected null", id, string(raw))
	}
	return nil
}

func (w *world) jsonTaskSpecEmpty(id string) error {
	specs, err := w.jsonSpecs(id)
	if err != nil {
		return err
	}
	if len(specs) != 0 {
		return fmt.Errorf("JSON task %s spec = %#v, expected an empty array", id, specs)
	}
	return nil
}

func (w *world) jsonTaskSpecEntry(id, path string, scenarios int) error {
	specs, err := w.jsonSpecs(id)
	if err != nil {
		return err
	}
	for _, s := range specs {
		if s.Path != path {
			continue
		}
		if s.State != "present" {
			return fmt.Errorf("JSON task %s spec %s state = %q, expected \"present\"", id, path, s.State)
		}
		if s.Scenarios != scenarios {
			return fmt.Errorf("JSON task %s spec %s scenarios = %d, expected %d", id, path, s.Scenarios, scenarios)
		}
		return nil
	}
	return fmt.Errorf("JSON task %s has no spec entry for %s; got: %#v", id, path, specs)
}

type jsonSpec struct {
	Path      string `json:"path"`
	State     string `json:"state"`
	Scenarios int    `json:"scenarios"`
}

func (w *world) jsonSpecField(id string) (json.RawMessage, error) {
	data, err := w.jsonObjectForTask(id)
	if err != nil {
		return nil, err
	}
	raw, ok := data["spec"]
	if !ok {
		return nil, fmt.Errorf("JSON task %s is missing field \"spec\"; task: %s", id, string(mustMarshal(data)))
	}
	return raw, nil
}

func (w *world) jsonSpecs(id string) ([]jsonSpec, error) {
	raw, err := w.jsonSpecField(id)
	if err != nil {
		return nil, err
	}
	var specs []jsonSpec
	if err := json.Unmarshal(raw, &specs); err != nil {
		return nil, fmt.Errorf("JSON task %s spec is not an array of specs (%v); got: %s", id, err, string(raw))
	}
	return specs, nil
}
