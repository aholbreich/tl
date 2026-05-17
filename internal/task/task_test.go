package task

import (
	"testing"
	"time"
)

func TestMarshalUnmarshalRoundtrip(t *testing.T) {
	actor := "claude-code:main"
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	orig := &Task{
		ID:        "task-abc123",
		Title:     "Add login validation",
		Status:    "open",
		Priority:  "high",
		Type:      "feature",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "human",
		Assignee:  nil,
		DependsOn: []string{"task-def456"},
		Claim: Claim{
			Actor:       &actor,
			ClaimedAt:   &now,
			ExpiresAt:   &now,
			HeartbeatAt: &now,
		},
		Verify: Verify{
			Commands:         []string{"go test ./..."},
			EvidenceRequired: []string{"test output"},
		},
		Tags: []string{"frontend", "auth"},
		Body: "## Description\n\nValidate email format.\n",
	}

	data, err := orig.MarshalMarkdown()
	if err != nil {
		t.Fatalf("MarshalMarkdown: %v", err)
	}

	parsed, err := UnmarshalMarkdown(data)
	if err != nil {
		t.Fatalf("UnmarshalMarkdown: %v", err)
	}

	if parsed.ID != orig.ID {
		t.Errorf("ID: got %q, want %q", parsed.ID, orig.ID)
	}
	if parsed.Title != orig.Title {
		t.Errorf("Title: got %q, want %q", parsed.Title, orig.Title)
	}
	if parsed.Status != orig.Status {
		t.Errorf("Status: got %q, want %q", parsed.Status, orig.Status)
	}
	if parsed.Priority != orig.Priority {
		t.Errorf("Priority: got %q, want %q", parsed.Priority, orig.Priority)
	}
	if parsed.Type != orig.Type {
		t.Errorf("Type: got %q, want %q", parsed.Type, orig.Type)
	}
	if len(parsed.DependsOn) != 1 || parsed.DependsOn[0] != "task-def456" {
		t.Errorf("DependsOn: got %v, want [task-def456]", parsed.DependsOn)
	}
	if parsed.Claim.Actor == nil || *parsed.Claim.Actor != "claude-code:main" {
		t.Errorf("Claim.Actor: got %v, want claude-code:main", parsed.Claim.Actor)
	}
	if len(parsed.Tags) != 2 {
		t.Errorf("Tags: got %v, want 2 tags", parsed.Tags)
	}
	if !parsed.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt: got %v, want %v", parsed.CreatedAt, now)
	}
	if len(parsed.Verify.Commands) != 1 || parsed.Verify.Commands[0] != "go test ./..." {
		t.Errorf("Verify.Commands: got %v", parsed.Verify.Commands)
	}
}

func TestMarshalNoClaimNoVerify(t *testing.T) {
	now := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)
	task := &Task{
		ID:        "task-min",
		Title:     "Minimal task",
		Status:    "open",
		Priority:  "medium",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "human",
		DependsOn: []string{},
		Tags:      []string{},
	}

	data, err := task.MarshalMarkdown()
	if err != nil {
		t.Fatalf("MarshalMarkdown: %v", err)
	}

	parsed, err := UnmarshalMarkdown(data)
	if err != nil {
		t.Fatalf("UnmarshalMarkdown: %v", err)
	}

	if parsed.ID != "task-min" {
		t.Errorf("ID: got %q", parsed.ID)
	}
	if parsed.Claim.Actor != nil {
		t.Errorf("Claim should be nil for minimal task, got %v", parsed.Claim.Actor)
	}
}

func TestMarshalWithBody(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	task := &Task{
		ID:        "task-body",
		Title:     "Body test",
		Status:    "open",
		Priority:  "medium",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "human",
		DependsOn: []string{},
		Tags:      []string{},
		Body:      "## Description\n\nSome description.\n\n## Notes\n\nSome note.\n",
	}

	data, err := task.MarshalMarkdown()
	if err != nil {
		t.Fatalf("MarshalMarkdown: %v", err)
	}

	parsed, err := UnmarshalMarkdown(data)
	if err != nil {
		t.Fatalf("UnmarshalMarkdown: %v", err)
	}

	if parsed.Body != task.Body {
		t.Errorf("Body: got %q, want %q", parsed.Body, task.Body)
	}
}

func TestUnmarshalMissingFrontmatter(t *testing.T) {
	_, err := UnmarshalMarkdown([]byte("not a task file"))
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}
