package cmd

import (
	"os"
	"testing"
)

func TestResolveActorCLIFlagTakesPriority(t *testing.T) {
	os.Setenv("TL_ACTOR", "env-agent")
	defer os.Unsetenv("TL_ACTOR")

	got := ResolveActor("cli-agent")
	if got != "cli-agent" {
		t.Errorf("ResolveActor: got %q, want cli-agent", got)
	}
}

func TestResolveActorEnvTL_ACTOR(t *testing.T) {
	os.Setenv("TL_ACTOR", "pi:main")
	defer os.Unsetenv("TL_ACTOR")

	got := ResolveActor("")
	if got != "pi:main" {
		t.Errorf("ResolveActor: got %q, want pi:main", got)
	}
}

func TestResolveActorEnvACTOR_NAME(t *testing.T) {
	os.Setenv("ACTOR_NAME", "fallback")
	defer os.Unsetenv("ACTOR_NAME")

	got := ResolveActor("")
	if got != "fallback" {
		t.Errorf("ResolveActor: got %q, want fallback", got)
	}
}

func TestResolveActorEnvBEADS_ACTOR(t *testing.T) {
	os.Setenv("BEADS_ACTOR", "beads")
	defer os.Unsetenv("BEADS_ACTOR")

	got := ResolveActor("")
	if got != "beads" {
		t.Errorf("ResolveActor: got %q, want beads", got)
	}
}

func TestResolveActorTL_ACTORPrecedesOthers(t *testing.T) {
	os.Setenv("TL_ACTOR", "primary")
	os.Setenv("ACTOR_NAME", "secondary")
	os.Setenv("BEADS_ACTOR", "tertiary")
	defer os.Unsetenv("TL_ACTOR")
	defer os.Unsetenv("ACTOR_NAME")
	defer os.Unsetenv("BEADS_ACTOR")

	got := ResolveActor("")
	if got != "primary" {
		t.Errorf("ResolveActor: got %q, want primary", got)
	}
}

func TestResolveActorAutoDetect(t *testing.T) {
	// Override the auto-detect function for testing.
	orig := DetectedActor
	DetectedActor = func() string { return "test-agent" }
	defer func() { DetectedActor = orig }()

	got := ResolveActor("")
	if got != "test-agent" {
		t.Errorf("ResolveActor auto-detect: got %q, want test-agent", got)
	}
}
