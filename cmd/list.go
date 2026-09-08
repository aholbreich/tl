package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/aholbreich/tl/internal/store"
	"github.com/aholbreich/tl/internal/task"
)

func newListCmd() *cobra.Command {
	var asJSON bool
	var includeAll bool
	var claimedBy string
	var status string
	var mine bool
	var tag string
	var dashboard bool
	var taskType string
	var priority string
	c := &cobra.Command{
		Use:   "list",
		Short: "List tasks in the ledger",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("priority") {
				var err error
				priority, err = normalizePriority(priority)
				if err != nil {
					return err
				}
			}
			ledger, err := requireLedger()
			if err != nil {
				return err
			}
			tasks, err := store.List(ledger)
			if err != nil {
				return err
			}
			tasks = filterListTasks(tasks, includeAll, claimedBy, status, mine, tag, taskType, priority)
			sortTasks(tasks)

			specs := resolveSpecsFor(ledger, tasks)
			showSpec := anySpecs(specs)

			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(compactTasksJSON(tasks, specs))
			}

			if dashboard {
				_, err := fmt.Fprint(cmd.OutOrStdout(), renderDashboard(tasks))
				return err
			}

			var rendered bytes.Buffer
			tw := tabwriter.NewWriter(&rendered, 0, 0, 2, ' ', 0)
			header := "ID\tStatus\tPriority\tClaimed By\tTitle"
			if showSpec {
				header += "\tSpec"
			}
			fmt.Fprintln(tw, header)
			for _, t := range tasks {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s", t.ID, t.Status, t.Priority, listClaimActor(t), t.Title)
				if showSpec {
					fmt.Fprintf(tw, "\t%s", specCell(specs[t.ID]))
				}
				fmt.Fprintln(tw)
			}
			if err := tw.Flush(); err != nil {
				return err
			}
			out := rendered.String()
			if useColor := commandColorEnabled(cmd); useColor {
				out = colorListRows(out, tasks, includeAll)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), out)
			return err
		},
	}
	c.Flags().BoolVar(&asJSON, "json", false, "Emit JSON output (takes precedence over --dashboard)")
	c.Flags().BoolVar(&dashboard, "dashboard", false, "Emit a status-grouped Markdown dashboard")
	c.Flags().StringVarP(&taskType, "type", "t", "", "Only show tasks of this type")
	c.Flags().StringVarP(&priority, "priority", "p", "", "Only show tasks with this priority (l/low|m/medium|h/high)")
	c.Flags().BoolVarP(&includeAll, "all", "a", false, "Include closed tasks (done and cancelled)")
	c.Flags().StringVar(&claimedBy, "claimed-by", "", "Only show tasks claimed by this actor")
	c.Flags().StringVar(&status, "status", "", "Only show tasks with this status (overrides default closed hiding)")
	c.Flags().BoolVar(&mine, "mine", false, "Only show tasks claimed by the resolved actor")
	c.Flags().StringVar(&tag, "tag", "", "Only show tasks carrying this tag")
	return c
}

func filterListTasks(tasks []*task.Task, includeAll bool, claimedBy string, status string, mine bool, tag, taskType, priority string) []*task.Task {
	if mine {
		resolved := ResolveActor("")
		claimedBy = resolved
	}

	filtered := tasks[:0]
	for _, t := range tasks {
		// --status overrides the default closed-task hiding.
		if status != "" {
			if t.Status != status {
				continue
			}
		} else if !includeAll && isClosedListStatus(t.Status) {
			continue
		}

		if claimedBy != "" && taskClaimActor(t) != claimedBy {
			continue
		}
		if tag != "" && !taskHasTag(t, tag) {
			continue
		}
		if taskType != "" && effectiveTaskType(t) != taskType {
			continue
		}
		if priority != "" && t.Priority != priority {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}

func isClosedListStatus(status string) bool {
	return status == "done" || status == "cancelled"
}

func colorListRows(rendered string, tasks []*task.Task, dimClosed bool) string {
	lines := strings.Split(rendered, "\n")
	var b strings.Builder
	for i, line := range lines {
		if i > 0 && i-1 < len(tasks) && line != "" {
			t := tasks[i-1]
			line = colorListRow(line, t, dimClosed && isClosedListStatus(t.Status))
		}
		b.WriteString(line)
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func colorListRow(line string, t *task.Task, dim bool) string {
	priorityStart := listPriorityStart(line, t)
	if priorityStart < 0 {
		return colorClosedListLine(dim, line)
	}
	priorityEnd := priorityStart + len(t.Priority)
	prefix := line[:priorityStart]
	priority := line[priorityStart:priorityEnd]
	suffix := line[priorityEnd:]

	if dim {
		return colorDimCode() + prefix + colorListPriority(true, priority) + colorDimCode() + suffix + colorResetCode()
	}
	return prefix + colorListPriority(true, priority) + suffix
}

func listPriorityStart(line string, t *task.Task) int {
	statusStart := strings.Index(line, t.Status)
	if statusStart < 0 {
		return -1
	}
	searchStart := statusStart + len(t.Status)
	priorityOffset := strings.Index(line[searchStart:], t.Priority)
	if priorityOffset < 0 {
		return -1
	}
	return searchStart + priorityOffset
}

func listClaimActor(t *task.Task) string {
	actor := taskClaimActor(t)
	if actor == "" {
		return "-"
	}
	return actor
}

func taskClaimActor(t *task.Task) string {
	if t.Claim.Actor == nil {
		return ""
	}
	return *t.Claim.Actor
}

func taskHasTag(t *task.Task, tag string) bool {
	for _, tg := range t.Tags {
		if tg == tag {
			return true
		}
	}
	return false
}

// statusSortOrder maps each status to a numeric rank (lower = appears first).
var statusSortOrder = map[string]int{
	"pending_human": 0,
	"blocked":       1,
	"in_progress":   2,
	"open":          3,
	"done":          4,
	"cancelled":     5,
}

func statusRank(s string) int {
	if r, ok := statusSortOrder[s]; ok {
		return r
	}
	return 99
}

// prioritySortRank maps priorities to numeric order (lower = appears first).
func prioritySortRank(priority string) int {
	switch priority {
	case "high":
		return 0
	case "medium":
		return 1
	case "low":
		return 2
	default:
		return 99
	}
}

// sortTasks orders tasks by status, then priority, then creation date (oldest first).
func sortTasks(tasks []*task.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		a, b := tasks[i], tasks[j]
		if sa, sb := statusRank(a.Status), statusRank(b.Status); sa != sb {
			return sa < sb
		}
		if a.Status != b.Status {
			return a.Status < b.Status
		}
		if pa, pb := prioritySortRank(a.Priority), prioritySortRank(b.Priority); pa != pb {
			return pa < pb
		}
		if !a.CreatedAt.Equal(b.CreatedAt) {
			return a.CreatedAt.Before(b.CreatedAt)
		}
		return a.ID < b.ID
	})
}
