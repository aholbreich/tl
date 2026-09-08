package cmd

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/aholbreich/tl/internal/task"
)

// renderDashboard consumes the same filtered, sorted tasks as the list table.
// No timestamps, terminal styles or writes: redirecting it produces a stable
// wiki snapshot without introducing another source of task state.
func renderDashboard(tasks []*task.Task) string {
	var out strings.Builder
	out.WriteString("# Task dashboard\n")
	if len(tasks) == 0 {
		out.WriteString("\nNo tasks match.\n")
		return out.String()
	}
	lastStatus := ""
	for i, t := range tasks {
		if i == 0 || t.Status != lastStatus {
			fmt.Fprintf(&out, "\n## %s\n", dashboardStatusHeading(t.Status))
			lastStatus = t.Status
		}
		fmt.Fprintf(&out, "\n### %s: %s\n\n", dashboardText(t.ID), dashboardText(t.Title))
		fmt.Fprintf(&out, "- Status: %s\n- Priority: %s\n- Type: %s\n- Claimant: %s\n",
			dashboardText(t.Status), dashboardText(t.Priority), dashboardText(effectiveTaskType(t)), dashboardText(listClaimActor(t)))
		if description := dashboardSummary(task.ParseBody(t.Body).Description); description != "" {
			fmt.Fprintf(&out, "\nDescription: %s\n", dashboardText(description))
		}
		if len(t.References) != 0 {
			out.WriteString("\nReferences:\n\n")
			for _, ref := range t.References {
				// References may be paths, URLs, IDs or free text. Do not turn
				// arbitrary strings into potentially unsafe Markdown links.
				fmt.Fprintf(&out, "- %s\n", dashboardText(ref))
			}
		}
	}
	return out.String()
}

func effectiveTaskType(t *task.Task) string {
	if t.Type == "" {
		return "task"
	}
	return t.Type
}

func dashboardStatusHeading(status string) string {
	switch status {
	case "pending_human":
		return "Pending human"
	case "blocked":
		return "Blocked"
	case "in_progress":
		return "In progress"
	case "open":
		return "Open"
	case "done":
		return "Done"
	case "cancelled":
		return "Cancelled"
	default:
		return dashboardText(status)
	}
}

// Summaries are one line, capped at 240 Unicode code points including ellipsis.
func dashboardSummary(description string) string {
	runes := []rune(strings.Join(strings.Fields(description), " "))
	if len(runes) > 240 {
		return string(runes[:239]) + "…"
	}
	return string(runes)
}

// Keep user text from injecting headings, HTML or terminal control sequences.
func dashboardText(text string) string {
	runes := []rune(strings.Join(strings.Fields(text), " "))
	var out strings.Builder
	for _, r := range runes {
		if unicode.IsControl(r) {
			continue
		}
		switch r {
		case '&':
			out.WriteString("&amp;")
		case '<':
			out.WriteString("&lt;")
		case '>':
			out.WriteString("&gt;")
		case '\\', '`', '*', '_', '[', ']', '#', '|':
			out.WriteByte('\\')
			out.WriteRune(r)
		default:
			out.WriteRune(r)
		}
	}
	return out.String()
}
