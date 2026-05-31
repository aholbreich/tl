package bdd

import (
	"bytes"
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"strings"

	"github.com/aholbreich/tl/cmd"
)

// --- agents.feature support -----------------------------------------------

func initializeAgentsSteps(ctx *godog.ScenarioContext, w *world) {
	ctx.Step(`^the file "([^"]*)" exists with content "([^"]*)"$`, w.fileExistsWithContent)
	ctx.Step(`^the file "([^"]*)" exists with content:$`, w.fileExistsWithDocString)
	ctx.Step(`^the file "([^"]*)" still has content "([^"]*)"$`, w.fileStillHasContent)
	ctx.Step(`^the file "([^"]*)" contains "([^"]*)"$`, w.fileContains)
	ctx.Step(`^the file "([^"]*)" does not contain "([^"]*)"$`, w.fileDoesNotContain)
	ctx.Step(`^the file "([^"]*)" does not exist$`, w.fileDoesNotExist)
	ctx.Step(`^the output contains "([^"]*)"$`, w.outputContains)
	ctx.Step(`^the output does not contain "([^"]*)"$`, w.outputDoesNotContain)
	ctx.Step(`^the output contains a "([^"]*)" heading$`, w.outputContainsHeading)
	ctx.Step(`^the output describes the ready, claim, show, note, and close steps$`, w.outputDescribesWorkflowSteps)
	ctx.Step(`^the output formats task commands as Markdown code spans$`, w.outputFormatsCommandsAsMarkdownCodeSpans)
	ctx.Step(`^the output contains these snippets:$`, w.outputContainsSnippets)
	ctx.Step(`^the compact agents output is shorter than the full agents output$`, w.compactAgentsOutputIsShorter)
}

func (w *world) fileExistsWithContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func (w *world) fileExistsWithDocString(path string, content *godog.DocString) error {
	return os.WriteFile(path, []byte(content.Content), 0o644)
}

func (w *world) fileStillHasContent(path, expected string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if string(data) != expected {
		return fmt.Errorf("file %s content = %q, expected %q", path, string(data), expected)
	}
	return nil
}

func (w *world) fileContains(path, needle string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), needle) {
		return fmt.Errorf("file %s does not contain %q; got:\n%s", path, needle, string(data))
	}
	return nil
}

func (w *world) fileDoesNotContain(path, needle string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), needle) {
		return fmt.Errorf("file %s contains %q; got:\n%s", path, needle, string(data))
	}
	return nil
}

func (w *world) fileDoesNotExist(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s exists", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (w *world) outputDoesNotContain(needle string) error {
	combined := w.stdout.String() + w.stderr.String()
	if w.cmdErr != nil {
		combined += "\n" + w.cmdErr.Error()
	}
	if strings.Contains(combined, needle) {
		return fmt.Errorf("output contains %q; got:\n%s", needle, combined)
	}
	return nil
}

func (w *world) outputContainsHeading(heading string) error {
	needle := "## " + heading
	if !strings.Contains(w.stdout.String(), needle) {
		return fmt.Errorf("output does not contain heading %q; got:\n%s", needle, w.stdout.String())
	}
	return nil
}

func (w *world) outputDescribesWorkflowSteps() error {
	for _, command := range []string{"tl ready", "tl claim", "tl show", "tl note", "tl close"} {
		if !strings.Contains(w.stdout.String(), command) {
			return fmt.Errorf("output does not describe %s; got:\n%s", command, w.stdout.String())
		}
	}
	return nil
}

func (w *world) outputFormatsCommandsAsMarkdownCodeSpans() error {
	for _, command := range []string{
		"`tl ready --json`",
		"`tl ready --tag <role> --json`",
		"`tl show <task-id>`",
		"`tl history <task-id>`",
		"`tl create \"<title>\" -d \"<description>\" --ref <path-or-url>`",
		"`tl claim <task-id> --actor agent-name`",
		"`tl note <task-id> -m \"...\" --actor agent-name`",
		"`tl dep add <task-id> --on <task-id>`",
		"`tl dep remove <task-id> --on <task-id>`",
		"`tl stale`",
		"`tl close <task-id> --actor agent-name`",
		"`tl cancel <task-id> -m \"<reason>\" --actor agent-name`",
		"`tl block <task-id> -m \"<blocker>\" --actor agent-name`",
		"`tl unblock <task-id> --actor agent-name`",
		"`tl pending <task-id> --question \"...\" --actor agent-name`",
		"`tl resolve <task-id> --answer \"...\" --actor agent-name`",
		"`tl release <task-id> --actor agent-name`",
	} {
		if !strings.Contains(w.stdout.String(), command) {
			return fmt.Errorf("output does not contain Markdown code span %s; got:\n%s", command, w.stdout.String())
		}
	}
	return nil
}

func (w *world) compactAgentsOutputIsShorter() error {
	root := cmd.NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"agents"})
	if err := root.Execute(); err != nil {
		return err
	}
	if w.stdout.Len() >= buf.Len() {
		return fmt.Errorf("compact output length %d is not shorter than full output length %d", w.stdout.Len(), buf.Len())
	}
	return nil
}

func (w *world) outputContainsSnippets(table *godog.Table) error {
	for rowIdx, row := range table.Rows[1:] {
		if len(row.Cells) != 1 {
			return fmt.Errorf("snippet row %d has %d cells, expected 1", rowIdx+1, len(row.Cells))
		}
		snippet := strings.TrimSpace(row.Cells[0].Value)
		if !strings.Contains(w.stdout.String(), snippet) {
			return fmt.Errorf("output does not contain snippet %q; got:\n%s", snippet, w.stdout.String())
		}
	}
	return nil
}
