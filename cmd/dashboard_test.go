package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/aholbreich/tl/internal/task"
)

func TestDashboardSummary(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{" \n ", ""},
		{"First line.\n\nSecond\tline.", "First line. Second line."},
		{strings.Repeat("界", 240), strings.Repeat("界", 240)},
		{strings.Repeat("界", 241), strings.Repeat("界", 239) + "…"},
	} {
		got := dashboardSummary(tc.input)
		if got != tc.want || !utf8.ValidString(got) {
			t.Errorf("summary = %q, want %q", got, tc.want)
		}
	}
}

func TestDashboardEscapesUserText(t *testing.T) {
	input := "**bold** _italic_ [link](javascript:alert(1)) <script>&\n# Heading\x1b"
	want := `\*\*bold\*\* \_italic\_ \[link\](javascript:alert(1)) &lt;script&gt;&amp; \# Heading`
	if got := dashboardText(input); got != want {
		t.Fatalf("escaped text = %q, want %q", got, want)
	}
	tasks := []*task.Task{{
		ID: "task-api", Title: input, Status: "open", Priority: "high",
		Body:       "## Description\n\n" + input + "\n\n## Notes\n\nDo not include notes.\n",
		References: []string{input},
	}}
	got := renderDashboard(tasks)
	if strings.Count(got, want) != 3 || strings.Contains(got, "Do not include notes") {
		t.Fatalf("title, description and reference must be escaped, with notes omitted:\n%s", got)
	}
}

func TestDashboardCode(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"in_progress", "`in_progress`"},
		{"aho", "`aho`"},
		{"", "`-`"},
		{" \t\n", "`-`"},
		{"review\n\t team\x1b", "`review team`"},
		{"**lead** <script>&", "`**lead** <script>&`"},
		{"agent`name", "``agent`name``"},
		{"`aho`", "`` `aho` ``"},
		{"team``lead`", "``` team``lead` ```"},
	} {
		if got := dashboardCode(tc.input); got != tc.want {
			t.Errorf("dashboardCode(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestDashboardCompactMetadata(t *testing.T) {
	actor := "`aho`\nreview"
	tasks := []*task.Task{{
		ID: "task-api", Title: "Document API", Status: "in_progress", Priority: "high",
		Type: "research_question", Claim: task.Claim{Actor: &actor},
	}}
	got := renderDashboard(tasks)
	want := "`in_progress` · **high** · `research_question` · 👤 `` `aho` review ``\n"
	if !strings.Contains(got, want) {
		t.Fatalf("dashboard missing compact metadata %q:\n%s", want, got)
	}
	for _, label := range []string{"- Status:", "- Priority:", "- Type:", "- Claimant:"} {
		if strings.Contains(got, label) {
			t.Errorf("dashboard still contains metadata bullet %q", label)
		}
	}
}

func TestDashboardOrdering(t *testing.T) {
	now := time.Now()
	tasks := []*task.Task{
		{ID: "task-z", Status: "open", Priority: "high", CreatedAt: now},
		{ID: "task-b", Status: "open", Priority: "low", CreatedAt: now.Add(-time.Hour)},
		{ID: "task-a", Status: "open", Priority: "high", CreatedAt: now},
		{ID: "task-c", Status: "blocked", Priority: "low", CreatedAt: now},
		{ID: "task-d", Status: "open", Priority: "high", CreatedAt: now.Add(-time.Hour)},
		{ID: "task-e", Status: "pending_human", Priority: "medium", CreatedAt: now},
	}
	sortTasks(tasks)
	var ids []string
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	want := []string{"task-e", "task-c", "task-d", "task-a", "task-z", "task-b"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("order = %v, want %v", ids, want)
	}
	got := renderDashboard(tasks)
	if strings.Count(got, "## Open\n") != 1 {
		t.Fatalf("expected one Open group:\n%s", got)
	}
}

func TestListDashboardReadOnlyAndJSONCompatibility(t *testing.T) {
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	run := func(args ...string) (string, error) {
		root := NewRootCmd()
		var out bytes.Buffer
		root.SetOut(&out)
		root.SetErr(&bytes.Buffer{})
		root.SetArgs(args)
		err := root.Execute()
		return out.String(), err
	}
	mustRun := func(args ...string) string {
		t.Helper()
		out, err := run(args...)
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return out
	}
	mustRun("init")
	mustRun("create", "Document API", "--type", "feature", "-p", "high", "--tag", "api", "-d", "Describe endpoints.")
	mustRun("create", "Review CLI", "--type", "chore", "-p", "low")
	before := dashboardLedgerSnapshot(t)
	first := mustRun("list", "--dashboard")
	if next := mustRun("list", "--dashboard"); next != first {
		t.Fatalf("dashboard changed between reads:\n%s\n%s", first, next)
	}
	plainJSON := mustRun("list", "--json", "--type", "feature", "--priority", "h")
	dashboardJSON := mustRun("list", "--dashboard", "--json", "--type", "feature", "--priority", "high")
	if plainJSON != dashboardJSON {
		t.Fatalf("dashboard must preserve list JSON:\n%s\n%s", plainJSON, dashboardJSON)
	}
	var tasks []task.Task
	if err := json.Unmarshal([]byte(plainJSON), &tasks); err != nil || len(tasks) != 1 || tasks[0].Title != "Document API" {
		t.Fatalf("filtered JSON = %s, error = %v", plainJSON, err)
	}
	if out := mustRun("list", "--type", "feature", "--priority", "h"); !strings.Contains(out, "Document API") || strings.Contains(out, "Review CLI") {
		t.Fatalf("table filtering failed: %s", out)
	}
	if out := mustRun("list", "--dashboard", "--type", "bug"); out != "# Task dashboard\n\nNo tasks match.\n" {
		t.Fatalf("empty filtered dashboard = %q", out)
	}
	if after := dashboardLedgerSnapshot(t); !reflect.DeepEqual(before, after) {
		t.Fatal("list changed ledger files")
	}
}

func dashboardLedgerSnapshot(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(".tl", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		files[path] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestListRejectsInvalidPriority(t *testing.T) {
	root := NewRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"list", "--dashboard", "--priority", "urgent"})
	err := root.Execute()
	var exitErr *ExitError
	if !errors.As(err, &exitErr) || exitErr.Code != 2 {
		t.Fatalf("invalid priority error = %v, want exit code 2", err)
	}
}
