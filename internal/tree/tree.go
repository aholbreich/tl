// Package tree renders the dependency graph the ledger already stores.
//
// A task's children are the tasks it depends on, so a parent sits above the
// slices it waits for. A root is a task that nothing else depends on.
//
// This package is a pure view over data tl already has: it adds no storage,
// no field, and no notion of progress.
package tree

import (
	"github.com/aholbreich/tl/internal/task"
)

// Node is one rendered position in the forest. The same task may appear as
// several nodes: a shared dependency is legitimate and is rendered under each
// parent rather than arbitrarily assigned to one.
type Node struct {
	Task     *task.Task
	Children []*Node
	// Cycle marks a node that closes a loop back onto one of its own
	// ancestors. Its children are not expanded, so rendering terminates.
	// Diagnosing cycles belongs to tl doctor; here they must merely survive.
	Cycle bool
}

// Build returns the forest for tasks, which the caller has already ordered.
//
// With rootID empty every root is returned. A task reachable from no root at
// all — which happens when a cycle has no entry point — is appended as its own
// root, so no visible task can silently vanish from the output.
//
// With rootID set only that subtree is returned, and the task is rendered even
// when the closed filter would otherwise hide it: naming a task explicitly is a
// stronger signal than the default.
func Build(tasks []*task.Task, rootID string, includeClosed bool) []*Node {
	byID := make(map[string]*task.Task, len(tasks))
	for _, t := range tasks {
		byID[t.ID] = t
	}

	visible := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		if includeClosed || !isClosed(t.Status) {
			visible[t.ID] = true
		}
	}

	if rootID != "" {
		root, ok := byID[rootID]
		if !ok {
			return nil
		}
		return []*Node{expand(root, byID, visible, map[string]bool{})}
	}

	// A root is a task no *visible* task depends on. Computing this over the
	// visible set rather than all tasks means an open task whose only parent
	// is closed surfaces as a root instead of disappearing with it.
	depended := make(map[string]bool)
	for _, t := range tasks {
		if !visible[t.ID] {
			continue
		}
		for _, dep := range t.DependsOn {
			if visible[dep] {
				depended[dep] = true
			}
		}
	}

	rendered := make(map[string]bool)
	var forest []*Node
	for _, t := range tasks {
		if !visible[t.ID] || depended[t.ID] {
			continue
		}
		n := expand(t, byID, visible, map[string]bool{})
		markRendered(n, rendered)
		forest = append(forest, n)
	}

	// Whatever no root reached — every task in a closed loop, for instance.
	for _, t := range tasks {
		if !visible[t.ID] || rendered[t.ID] {
			continue
		}
		n := expand(t, byID, visible, map[string]bool{})
		markRendered(n, rendered)
		forest = append(forest, n)
	}
	return forest
}

// expand walks depth-first. path holds the ancestors of the current node, so a
// repeat within it is a cycle rather than the diamond of a shared dependency.
func expand(t *task.Task, byID map[string]*task.Task, visible, path map[string]bool) *Node {
	if path[t.ID] {
		return &Node{Task: t, Cycle: true}
	}
	path[t.ID] = true
	defer delete(path, t.ID)

	n := &Node{Task: t}
	for _, depID := range t.DependsOn {
		dep, ok := byID[depID]
		if !ok || !visible[depID] {
			// A dangling dependency is tl doctor's to report, not the
			// tree's to render.
			continue
		}
		n.Children = append(n.Children, expand(dep, byID, visible, path))
	}
	return n
}

func markRendered(n *Node, rendered map[string]bool) {
	rendered[n.Task.ID] = true
	for _, c := range n.Children {
		markRendered(c, rendered)
	}
}

func isClosed(status string) bool {
	return status == "done" || status == "cancelled"
}
