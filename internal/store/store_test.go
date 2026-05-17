package store

import (
	"os"
	"testing"
)

func TestPriorityRank(t *testing.T) {
	tests := []struct {
		priority string
		want     int
	}{
		{"high", 0},
		{"medium", 1},
		{"low", 2},
		{"unknown", 3},
		{"", 3},
	}
	for _, tt := range tests {
		got := priorityRank(tt.priority)
		if got != tt.want {
			t.Errorf("priorityRank(%q) = %d, want %d", tt.priority, got, tt.want)
		}
	}
}

func TestNormalizeID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"task-abc123", "task-abc123"},
		{"abc123", "task-abc123"},
		{"task-", "task-"},
		{"", "task-"},
	}
	for _, tt := range tests {
		got := NormalizeID(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNewIDFormat(t *testing.T) {
	dir := t.TempDir()
	ledger := dir + "/.taskledger"
	if err := os.MkdirAll(ledger+"/tasks", 0755); err != nil {
		t.Fatalf("create tasks dir: %v", err)
	}

	id, err := NewID(ledger)
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	if len(id) != len("task-xxx") {
		t.Errorf("id length: got %d, want %d (%q)", len(id), len("task-xxx"), id)
	}
	if id[:5] != "task-" {
		t.Errorf("id prefix: got %q, want task-", id[:5])
	}
}
