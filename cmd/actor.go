package cmd

import "os"

// ActorEnvChain is the environment variable fallback chain for automatic actor
// resolution. Earlier entries take precedence.
var ActorEnvChain = []string{"TL_ACTOR", "ACTOR_NAME", "BEADS_ACTOR"}

// DefaultDetectActor is the default auto-detection function, exported so
// tests can restore DetectedActor after overriding it.
var DefaultDetectActor = defaultDetectActor

// DetectedActor is called when no explicit actor is provided via CLI flag or
// environment variable. Tests may override this to simulate agent detection.
var DetectedActor = DefaultDetectActor

// ResolveActor returns the actor identity using this priority:
//  1. CLI --actor flag (when non-empty)
//  2. Environment variables in ActorEnvChain order
//  3. Auto-detection via DetectedActor
func ResolveActor(cliFlag string) string {
	if cliFlag != "" {
		return cliFlag
	}
	for _, env := range ActorEnvChain {
		if v := os.Getenv(env); v != "" {
			return v
		}
	}
	return DetectedActor()
}

// agentMarker identifies a coding harness from its environment.
type agentMarker struct {
	env   string
	actor string
	// value, when set, must match exactly rather than merely being present.
	// A harness that exports a boolean flag says nothing by defining it as
	// "false".
	value string
}

// agentMarkers is ordered by precedence: the first match wins.
//
// Only markers a harness sets reliably belong here, because a false positive
// silently attributes one agent's claims to another. Two deliberate
// omissions: Cursor, whose VSCODE_* variables are set by plain VS Code as
// well, and GitHub Copilot, whose GITHUB_COPILOT_TOKEN is an auth credential
// that outlives any session rather than a marker of one.
var agentMarkers = []agentMarker{
	{env: "PI_CODING_AGENT", actor: "pi", value: "true"},
	{env: "CLAUDE_CODE_SESSION_ID", actor: "claude-code"},
	{env: "AIDER_MODEL", actor: "aider"},
	{env: "CODEIUM_API_KEY", actor: "windsurf"},
	{env: "PI_AGENT_ID", actor: "pi"},
}

// ActorDetectionEnv returns every environment variable defaultDetectActor
// consults. Test harnesses clear these before a scenario, so that a marker
// that happens to be set on the machine running the suite cannot decide the
// result — detection tests would otherwise pass or fail based on which
// harness the developer was using at the time.
func ActorDetectionEnv() []string {
	out := make([]string, 0, len(agentMarkers))
	for _, m := range agentMarkers {
		out = append(out, m.env)
	}
	return out
}

func defaultDetectActor() string {
	for _, m := range agentMarkers {
		v := os.Getenv(m.env)
		if v == "" {
			continue
		}
		if m.value != "" && v != m.value {
			continue
		}
		return m.actor
	}
	if _, err := os.Stat(".codex"); err == nil {
		return "codex"
	}
	// Fallback: use hostname.
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	return "unknown"
}
