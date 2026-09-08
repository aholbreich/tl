package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/aholbreich/tl/internal/spec"
	"github.com/aholbreich/tl/internal/store"
	"github.com/aholbreich/tl/internal/task"
	"github.com/aholbreich/tl/internal/tree"
)

func newTreeCmd() *cobra.Command {
	var asJSON bool
	var includeAll bool
	var specStatus bool
	c := &cobra.Command{
		Use:               "tree [TASK_ID]",
		Short:             "Render the dependency graph",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: completeTaskIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ledger, err := requireLedger()
			if err != nil {
				return err
			}
			tasks, err := store.List(ledger)
			if err != nil {
				return err
			}
			sortTasks(tasks)

			rootID := ""
			if len(args) == 1 {
				rootID = store.NormalizeID(args[0])
				if _, err := store.Read(ledger, rootID); err != nil {
					return NewExitError(3, "task %s not found", args[0])
				}
			}

			forest := tree.Build(tasks, rootID, includeAll)
			specs := resolveSpecsFor(ledger, tasks, specStatus)

			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(treeForestJSON(forest))
			}

			if len(forest) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "No tasks.")
				return err
			}
			return renderForest(cmd, forest, includeAll, specs)
		},
	}
	c.Flags().BoolVar(&asJSON, "json", false, "Emit JSON output")
	c.Flags().BoolVarP(&includeAll, "all", "a", false, "Include closed tasks (done and cancelled)")
	c.Flags().BoolVar(&specStatus, "spec-status", false, specStatusFlagUsage)
	return c
}

// renderForest lays the tree out with a tabwriter and only then applies color.
// Doing it the other way round would feed ANSI escapes into the width
// calculation and misalign every row, which is why list does the same.
func renderForest(cmd *cobra.Command, forest []*tree.Node, dimClosed bool, specs map[string][]spec.Spec) error {
	var rendered bytes.Buffer
	tw := tabwriter.NewWriter(&rendered, 0, 0, 2, ' ', 0)
	var order []*tree.Node
	for _, root := range forest {
		writeNode(tw, root, "", "", &order, specs)
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	out := rendered.String()
	if commandColorEnabled(cmd) {
		out = colorTreeRows(out, order, dimClosed)
	}
	_, err := fmt.Fprint(cmd.OutOrStdout(), out)
	return err
}

// writeNode emits one row and recurses. prefix is what precedes this row's
// branch glyph; childPrefix is what every descendant row carries instead.
func writeNode(tw *tabwriter.Writer, n *tree.Node, prefix, childPrefix string, order *[]*tree.Node, specs map[string][]spec.Spec) {
	suffix := ""
	if n.Cycle {
		suffix = "  (cycle)"
	}
	fmt.Fprintf(tw, "%s%s\t%s\t%s\t%s%s", prefix, n.Task.ID, n.Task.Title, n.Task.Status, n.Task.Priority, suffix)
	if specs != nil {
		fmt.Fprintf(tw, "\t%s", specCell(specs[n.Task.ID]))
	}
	fmt.Fprintln(tw)
	*order = append(*order, n)

	for i, child := range n.Children {
		last := i == len(n.Children)-1
		branch, cont := "├─ ", "│  "
		if last {
			branch, cont = "└─ ", "   "
		}
		writeNode(tw, child, childPrefix+branch, childPrefix+cont, order, specs)
	}
}

func colorTreeRows(rendered string, order []*tree.Node, dimClosed bool) string {
	lines := strings.Split(rendered, "\n")
	var b strings.Builder
	for i, line := range lines {
		if i < len(order) && line != "" {
			line = colorTreeRow(line, order[i], dimClosed && isClosedListStatus(order[i].Task.Status))
		}
		b.WriteString(line)
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// colorTreeRow mirrors colorListRow so the two views cannot drift apart.
func colorTreeRow(line string, n *tree.Node, dim bool) string {
	start := treePriorityStart(line, n.Task)
	if start < 0 {
		return colorClosedListLine(dim, line)
	}
	end := start + len(n.Task.Priority)
	prefix, priority, suffix := line[:start], line[start:end], line[end:]

	if dim {
		return colorDimCode() + prefix + colorListPriority(true, priority) + colorDimCode() + suffix + colorResetCode()
	}
	return prefix + colorListPriority(true, priority) + suffix
}

// treePriorityStart locates the priority cell. Unlike the list table, a tree
// row carries the title before the status, so the search is anchored on the
// exact title first — otherwise a task called "open the door" would have its
// own title mistaken for the status column.
func treePriorityStart(line string, t *task.Task) int {
	titleAt := strings.Index(line, t.Title)
	if titleAt < 0 {
		return -1
	}
	searchFrom := titleAt + len(t.Title)
	statusAt := strings.Index(line[searchFrom:], t.Status)
	if statusAt < 0 {
		return -1
	}
	searchFrom += statusAt + len(t.Status)
	priorityAt := strings.Index(line[searchFrom:], t.Priority)
	if priorityAt < 0 {
		return -1
	}
	return searchFrom + priorityAt
}

// treeNodeJSON nests children so a consumer does not re-derive the graph.
// It is deliberately leaner than the task JSON of list and show: a shared
// dependency appears under every parent, and repeating whole task bodies at
// each position would make the document large and confusing rather than
// useful. Consumers that need the full task fetch it by id.
type treeNodeJSON struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	Status   string         `json:"status"`
	Priority string         `json:"priority"`
	Cycle    bool           `json:"cycle,omitempty"`
	Children []treeNodeJSON `json:"children"`
}

func treeForestJSON(forest []*tree.Node) []treeNodeJSON {
	out := make([]treeNodeJSON, 0, len(forest))
	for _, n := range forest {
		out = append(out, treeNodeJSON1(n))
	}
	return out
}

func treeNodeJSON1(n *tree.Node) treeNodeJSON {
	children := make([]treeNodeJSON, 0, len(n.Children))
	for _, c := range n.Children {
		children = append(children, treeNodeJSON1(c))
	}
	return treeNodeJSON{
		ID:       n.Task.ID,
		Title:    n.Task.Title,
		Status:   n.Task.Status,
		Priority: n.Task.Priority,
		Cycle:    n.Cycle,
		Children: children,
	}
}
