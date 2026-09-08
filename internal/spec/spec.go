// Package spec resolves task references that point at specification files.
//
// Gherkin `.feature` files are the first and currently the only resolver.
// The exported names deliberately talk about references and specs rather
// than Gherkin, so a second format can be added without a redesign.
//
// Resolution is file-granular by decision 0002: a reference names a file,
// and this package reports what can be observed about that file. It never
// interprets what it finds — no tag is given meaning here, and nothing in
// this package concludes that work is done.
package spec

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// State is what could be observed about a referenced specification file.
type State string

const (
	// StatePresent means the file was read successfully.
	StatePresent State = "present"
	// StateMissing means the reference points at a file that is not there.
	StateMissing State = "missing"
	// StateUnknown means the file exists but could not be read. A listing
	// degrades to this rather than failing; tl doctor owns the complaining.
	StateUnknown State = "unknown"
)

// Spec is one resolved specification reference.
type Spec struct {
	Path  string `json:"path"`
	State State  `json:"state"`
	// Scenarios counts scenario blocks, so a Scenario Outline counts once
	// however many Examples rows it carries. Zero unless State is present.
	Scenarios int `json:"scenarios"`
}

// specSuffix is the one entry in the format table today. A second resolver
// adds a suffix here and a counter below; nothing else changes.
const specSuffix = ".feature"

var urlSchemeRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)

// IsRef reports whether ref names a specification file. It is a string rule
// and touches no disk, so it is safe on the default path.
func IsRef(ref string) bool {
	if urlSchemeRE.MatchString(ref) {
		return false // a URL is a link, not a file this repository owns
	}
	return strings.HasSuffix(ref, specSuffix)
}

// Refs returns the subset of refs that name specification files, preserving
// order.
func Refs(refs []string) []string {
	var out []string
	for _, ref := range refs {
		if IsRef(ref) {
			out = append(out, ref)
		}
	}
	return out
}

// Resolve reads every specification reference in refs, relative to repoRoot.
// The result is non-nil so callers can distinguish "resolved, none found"
// from "not resolved at all".
func Resolve(repoRoot string, refs []string) []Spec {
	out := make([]Spec, 0, len(refs))
	for _, ref := range Refs(refs) {
		out = append(out, resolveOne(repoRoot, ref))
	}
	return out
}

func resolveOne(repoRoot, ref string) Spec {
	f, err := os.Open(filepath.Join(repoRoot, ref))
	if err != nil {
		if os.IsNotExist(err) {
			return Spec{Path: ref, State: StateMissing}
		}
		return Spec{Path: ref, State: StateUnknown}
	}
	defer f.Close()

	n, err := countScenarios(f)
	if err != nil {
		return Spec{Path: ref, State: StateUnknown}
	}
	return Spec{Path: ref, State: StatePresent, Scenarios: n}
}

// scenarioKeywords are the Gherkin keywords that open a scenario block.
// "Examples:" is deliberately absent: it is the data table belonging to an
// outline, not a scenario of its own.
var scenarioKeywords = []string{
	"Scenario:",
	"Scenario Outline:",
	"Scenario Template:",
	"Example:",
}

func countScenarios(f *os.File) (int, error) {
	count := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		for _, kw := range scenarioKeywords {
			if strings.HasPrefix(line, kw) {
				count++
				break
			}
		}
	}
	return count, sc.Err()
}
