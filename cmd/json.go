package cmd

import (
	"time"

	"github.com/aholbreich/tl/internal/spec"
	"github.com/aholbreich/tl/internal/task"
)

type compactTaskJSON struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Status      string            `json:"status"`
	Priority    string            `json:"priority"`
	Type        string            `json:"type,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedBy   string            `json:"created_by"`
	Assignee    *string           `json:"assignee"`
	DependsOn   []string          `json:"depends_on"`
	References  []string          `json:"references"`
	Spec        []spec.Spec       `json:"spec"`
	Claim       task.Claim        `json:"claim"`
	Pending     *task.Pending     `json:"pending,omitempty"`
	Tags        []string          `json:"tags"`
	Description string            `json:"description,omitempty"`
	Notes       []task.Note       `json:"notes,omitempty"`
	Sections    map[string]string `json:"sections,omitempty"`
}

// compactTasksJSON renders tasks for bulk output. Every task carries a spec
// entry: an array, empty when it references no specification. There is no
// null case, so a consumer can index it without a presence check.
func compactTasksJSON(tasks []*task.Task, specs map[string][]spec.Spec) []compactTaskJSON {
	out := make([]compactTaskJSON, 0, len(tasks))
	for _, t := range tasks {
		parsed := task.ParseBody(t.Body)
		out = append(out, compactTaskJSON{
			ID:          t.ID,
			Title:       t.Title,
			Status:      t.Status,
			Priority:    t.Priority,
			Type:        t.Type,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			CreatedBy:   t.CreatedBy,
			Assignee:    t.Assignee,
			DependsOn:   t.DependsOn,
			References:  compactReferences(t.References),
			Spec:        specs[t.ID],
			Claim:       t.Claim,
			Pending:     t.Pending,
			Tags:        t.Tags,
			Description: parsed.Description,
			Notes:       parsed.Notes,
			Sections:    parsed.Sections,
		})
	}
	return out
}

// compactReferences guarantees an empty array rather than null, matching the
// shape `tl show --json` emits for a task with no references.
func compactReferences(refs []string) []string {
	if refs == nil {
		return []string{}
	}
	return refs
}
