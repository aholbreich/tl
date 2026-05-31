package cmd

import (
	"strings"
	"testing"
)

func TestMergeAgentsBlockAppendsToExistingContent(t *testing.T) {
	got := mergeAgentsBlock("# Project", agentsSnippet)
	for _, want := range []string{"# Project\n\n", agentsBeginMarker, "## tl workflow", agentsEndMarker} {
		if !strings.Contains(got, want) {
			t.Fatalf("mergeAgentsBlock() missing %q; got:\n%s", want, got)
		}
	}
}

func TestMergeAgentsBlockReplacesManagedBlockIdempotently(t *testing.T) {
	input := "# Project\n\n" + agentsBeginMarker + "\nold workflow\n" + agentsEndMarker + "\n"

	once := mergeAgentsBlock(input, agentsSnippet)
	twice := mergeAgentsBlock(once, agentsSnippet)

	if once != twice {
		t.Fatalf("mergeAgentsBlock should be idempotent\nonce:\n%s\ntwice:\n%s", once, twice)
	}
	if strings.Contains(once, "old workflow") {
		t.Fatalf("mergeAgentsBlock did not replace old managed block; got:\n%s", once)
	}
}

func TestMergeAgentsBlockReplacesLegacyManagedBlock(t *testing.T) {
	input := "# Project\n\n" + legacyAgentsBeginMarker + "\nold workflow\n" + legacyAgentsEndMarker + "\n"

	got := mergeAgentsBlock(input, agentsSnippet)

	if strings.Contains(got, "old workflow") {
		t.Fatalf("mergeAgentsBlock did not replace legacy managed block; got:\n%s", got)
	}
	if !strings.Contains(got, agentsBeginMarker) || !strings.Contains(got, agentsEndMarker) {
		t.Fatalf("mergeAgentsBlock did not write current markers; got:\n%s", got)
	}
}
