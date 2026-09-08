package bdd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// --- tree.feature support -------------------------------------------------

func initializeTreeSteps(ctx *godog.ScenarioContext, w *world) {
	ctx.Step(`^the output shows "([^"]*)" as a child of "([^"]*)"$`, w.treeShowsChildOf)
	ctx.Step(`^the output shows "([^"]*)" at depth (\d+)$`, w.treeShowsAtDepth)
	ctx.Step(`^the output shows "([^"]*)" (\d+) times$`, w.treeShowsNTimes)
	ctx.Step(`^the output marks "([^"]*)" as a cycle$`, w.treeMarksCycle)
	ctx.Step(`^the tree row for "([^"]*)" contains "([^"]*)"$`, w.treeRowContains)
	ctx.Step(`^the JSON tree root "([^"]*)" has a child "([^"]*)"$`, w.jsonTreeRootHasChild)
	ctx.Step(`^the JSON tree root "([^"]*)" has no children$`, w.jsonTreeRootHasNoChildren)
	ctx.Step(`^the JSON tree root "([^"]*)" has a spec entry for "([^"]*)"$`, w.jsonTreeRootHasSpec)
	ctx.Step(`^the JSON tree child "([^"]*)" has a spec entry for "([^"]*)"$`, w.jsonTreeChildHasSpec)
	ctx.Step(`^the JSON tree root "([^"]*)" has an empty spec array$`, w.jsonTreeRootSpecEmpty)
}

// treeDepth derives nesting from the branch glyphs preceding the identifier,
// so the assertion reads the drawn tree rather than trusting the writer.
func treeDepth(line string) int {
	prefix := line[:strings.Index(line, "task-")]
	return strings.Count(prefix, "├─ ") + strings.Count(prefix, "└─ ") +
		strings.Count(prefix, "│  ") + strings.Count(prefix, "   ")
}

func (w *world) treeLinesFor(id string) []string {
	var out []string
	for _, line := range strings.Split(w.stdout.String(), "\n") {
		if strings.Contains(line, id) {
			out = append(out, line)
		}
	}
	return out
}

func (w *world) treeLineFor(id string) (string, error) {
	lines := w.treeLinesFor(id)
	if len(lines) == 0 {
		return "", fmt.Errorf("output has no row for %s; got:\n%s", id, w.stdout.String())
	}
	return lines[0], nil
}

func (w *world) treeShowsAtDepth(id string, depth int) error {
	for _, line := range w.treeLinesFor(id) {
		if treeDepth(line) == depth {
			return nil
		}
	}
	return fmt.Errorf("no row for %s at depth %d; got:\n%s", id, depth, w.stdout.String())
}

func (w *world) treeShowsChildOf(child, parent string) error {
	parentLine, err := w.treeLineFor(parent)
	if err != nil {
		return err
	}
	childLine, err := w.treeLineFor(child)
	if err != nil {
		return err
	}
	if treeDepth(childLine) != treeDepth(parentLine)+1 {
		return fmt.Errorf("%s is not one level below %s; got:\n%s", child, parent, w.stdout.String())
	}
	if !strings.Contains(w.stdout.String(), parentLine+"\n") {
		return fmt.Errorf("expected %s to precede %s", parent, child)
	}
	return nil
}

func (w *world) treeShowsNTimes(id string, want int) error {
	got := len(w.treeLinesFor(id))
	if got != want {
		return fmt.Errorf("expected %s %d times, found %d; got:\n%s", id, want, got, w.stdout.String())
	}
	return nil
}

func (w *world) treeRowContains(id, want string) error {
	line, err := w.treeLineFor(id)
	if err != nil {
		return err
	}
	if !strings.Contains(line, want) {
		return fmt.Errorf("tree row for %s does not contain %q; row: %s", id, want, line)
	}
	return nil
}

func (w *world) treeMarksCycle(id string) error {
	for _, line := range w.treeLinesFor(id) {
		if strings.Contains(line, "(cycle)") {
			return nil
		}
	}
	return fmt.Errorf("no row for %s marked as a cycle; got:\n%s", id, w.stdout.String())
}

// --- JSON -----------------------------------------------------------------

type jsonTreeNode struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	Status   string         `json:"status"`
	Cycle    bool           `json:"cycle"`
	Spec     []jsonSpec     `json:"spec"`
	Children []jsonTreeNode `json:"children"`
}

// findTreeNode walks the whole forest, so a spec assertion works at any depth
// rather than only at a root.
func findTreeNode(nodes []jsonTreeNode, id string) (jsonTreeNode, bool) {
	for _, n := range nodes {
		if n.ID == id {
			return n, true
		}
		if found, ok := findTreeNode(n.Children, id); ok {
			return found, true
		}
	}
	return jsonTreeNode{}, false
}

func (w *world) jsonTreeNodeSpecs(id string) ([]jsonSpec, error) {
	var forest []jsonTreeNode
	if err := json.Unmarshal(w.stdout.Bytes(), &forest); err != nil {
		return nil, fmt.Errorf("stdout is not a JSON tree (%v); got: %s", err, w.stdout.String())
	}
	n, ok := findTreeNode(forest, id)
	if !ok {
		return nil, fmt.Errorf("JSON tree has no node %s; got: %s", id, w.stdout.String())
	}
	if n.Spec == nil {
		return nil, fmt.Errorf("JSON node %s has a null spec; expected an array", id)
	}
	return n.Spec, nil
}

func (w *world) jsonTreeRootHasSpec(id, path string) error {
	return w.jsonTreeNodeHasSpec(id, path)
}

func (w *world) jsonTreeChildHasSpec(id, path string) error {
	return w.jsonTreeNodeHasSpec(id, path)
}

func (w *world) jsonTreeNodeHasSpec(id, path string) error {
	specs, err := w.jsonTreeNodeSpecs(id)
	if err != nil {
		return err
	}
	for _, sp := range specs {
		if sp.Path == path {
			return nil
		}
	}
	return fmt.Errorf("JSON node %s has no spec entry for %s; got %#v", id, path, specs)
}

func (w *world) jsonTreeRootSpecEmpty(id string) error {
	specs, err := w.jsonTreeNodeSpecs(id)
	if err != nil {
		return err
	}
	if len(specs) != 0 {
		return fmt.Errorf("JSON node %s spec = %#v, expected empty", id, specs)
	}
	return nil
}

func (w *world) jsonTreeRoot(id string) (jsonTreeNode, error) {
	var forest []jsonTreeNode
	if err := json.Unmarshal(w.stdout.Bytes(), &forest); err != nil {
		return jsonTreeNode{}, fmt.Errorf("stdout is not a JSON tree (%v); got: %s", err, w.stdout.String())
	}
	for _, n := range forest {
		if n.ID == id {
			return n, nil
		}
	}
	return jsonTreeNode{}, fmt.Errorf("JSON tree has no root %s; got: %s", id, w.stdout.String())
}

func (w *world) jsonTreeRootHasChild(rootID, childID string) error {
	root, err := w.jsonTreeRoot(rootID)
	if err != nil {
		return err
	}
	for _, c := range root.Children {
		if c.ID == childID {
			return nil
		}
	}
	return fmt.Errorf("JSON root %s has no child %s; got: %s", rootID, childID, w.stdout.String())
}

func (w *world) jsonTreeRootHasNoChildren(rootID string) error {
	root, err := w.jsonTreeRoot(rootID)
	if err != nil {
		return err
	}
	if root.Children == nil {
		return fmt.Errorf("JSON root %s children is null, expected an empty array", rootID)
	}
	if len(root.Children) != 0 {
		return fmt.Errorf("JSON root %s has %d children, expected none", rootID, len(root.Children))
	}
	return nil
}
