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

// clearActorDetectionEnv isolates a detection test from the machine running
// it. The suite is routinely executed inside a coding harness that exports
// one of these, which would otherwise decide every result below.
func clearActorDetectionEnv(t *testing.T) {
	t.Helper()
	for _, key := range ActorDetectionEnv() {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

func TestDetectActorFromMarkers(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"aider", map[string]string{"AIDER_MODEL": "gpt-4o"}, "aider"},
		{"windsurf", map[string]string{"CODEIUM_API_KEY": "sk-example"}, "windsurf"},
		{"pi via coding agent flag", map[string]string{"PI_CODING_AGENT": "true"}, "pi"},
		{"pi via agent id", map[string]string{"PI_AGENT_ID": "agent-7"}, "pi"},
		{"claude code", map[string]string{"CLAUDE_CODE_SESSION_ID": "0f9c1a2b"}, "claude-code"},
		{
			"claude code precedes aider",
			map[string]string{"CLAUDE_CODE_SESSION_ID": "0f9c1a2b", "AIDER_MODEL": "gpt-4o"},
			"claude-code",
		},
		{
			"pi flag precedes claude code",
			map[string]string{"PI_CODING_AGENT": "true", "CLAUDE_CODE_SESSION_ID": "0f9c1a2b"},
			"pi",
		},
		{
			"false pi flag falls through",
			map[string]string{"PI_CODING_AGENT": "false", "AIDER_MODEL": "gpt-4o"},
			"aider",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearActorDetectionEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			if got := defaultDetectActor(); got != tc.want {
				t.Errorf("defaultDetectActor: got %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDetectActorFallsBackToHostname pins the behaviour that makes an
// unrecognised harness usable at all: without it every such caller would
// share one identity and collide on claims.
func TestDetectActorFallsBackToHostname(t *testing.T) {
	clearActorDetectionEnv(t)
	t.Chdir(t.TempDir()) // away from any .codex marker in the repo

	host, err := os.Hostname()
	if err != nil || host == "" {
		t.Skip("hostname unavailable")
	}
	if got := defaultDetectActor(); got != host {
		t.Errorf("defaultDetectActor: got %q, want hostname %q", got, host)
	}
}

func TestActorDetectionEnvCoversEveryMarker(t *testing.T) {
	if got, want := len(ActorDetectionEnv()), len(agentMarkers); got != want {
		t.Errorf("ActorDetectionEnv lists %d vars for %d markers; a marker added without it would leak into tests", got, want)
	}
}
