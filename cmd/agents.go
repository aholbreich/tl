package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed agents_snippet.md
var agentsSnippet string

//go:embed agents_snippet_compact.md
var agentsSnippetCompact string

const (
	agentsBeginMarker       = "<!-- BEGIN TL WORKFLOW -->"
	agentsEndMarker         = "<!-- END TL WORKFLOW -->"
	legacyAgentsBeginMarker = "<!-- BEGIN " + "TASK" + "LEDGER WORKFLOW -->"
	legacyAgentsEndMarker   = "<!-- END " + "TASK" + "LEDGER WORKFLOW -->"
)

var agentInstructionFiles = []string{
	"AGENTS.md",
	"CLAUDE.md",
	"GEMINI_RULES.md",
	".cursorrules",
	".aider-rules.md",
	".github/copilot-instructions.md",
}

func newAgentsCmd() *cobra.Command {
	var writeFiles bool
	var dryRun bool
	var compact bool
	var files []string
	c := &cobra.Command{
		Use:   "agents",
		Short: "Print recommended AGENTS.md instructions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRun && !writeFiles {
				return NewExitError(2, "--dry-run requires --write-files")
			}
			if writeFiles {
				return updateAgentInstructionFiles(cmd, dryRun, compact, files)
			}
			_, err := fmt.Fprint(cmd.OutOrStdout(), selectedAgentsSnippet(compact))
			return err
		},
	}
	c.Flags().BoolVar(&writeFiles, "write-files", false, "Write or refresh the tl workflow block in existing agent instruction files")
	c.Flags().BoolVar(&writeFiles, "update", false, "(deprecated: use --write-files)")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "Report which agent instruction files would be updated without modifying them")
	c.Flags().BoolVar(&compact, "compact", false, "Print or write a compact tl workflow guide for constrained context windows")
	c.Flags().StringArrayVar(&files, "file", nil, "Only consider this agent instruction file (repeatable; defaults to known files)")
	_ = c.Flags().MarkHidden("update")
	return c
}

type agentInstructionFilePlan struct {
	Path            string
	Info            os.FileInfo
	Missing         bool
	HasManagedBlock bool
	Content         []byte
}

func updateAgentInstructionFiles(cmd *cobra.Command, dryRun, compact bool, files []string) error {
	paths := agentInstructionPaths(files)
	plans, err := scanAgentInstructionFiles(paths)
	if err != nil {
		return err
	}
	if dryRun {
		for _, plan := range plans {
			fmt.Fprintln(cmd.OutOrStdout(), plan.DryRunMessage())
		}
		return nil
	}

	snippet := selectedAgentsSnippet(compact)
	updated := 0
	for _, plan := range plans {
		if plan.Missing {
			if len(files) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Skipped %s (file not found)\n", plan.Path)
			}
			continue
		}
		if err := os.WriteFile(plan.Path, []byte(mergeAgentsBlock(string(plan.Content), snippet)), plan.Info.Mode()); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", plan.Path)
		updated++
	}
	if updated == 0 && len(files) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No existing agent instruction files found")
	}
	return nil
}

func agentInstructionPaths(files []string) []string {
	if len(files) > 0 {
		return files
	}
	return agentInstructionFiles
}

func scanAgentInstructionFiles(paths []string) ([]agentInstructionFilePlan, error) {
	plans := make([]agentInstructionFilePlan, 0, len(paths))
	for _, path := range paths {
		plan := agentInstructionFilePlan{Path: path}
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				plan.Missing = true
				plans = append(plans, plan)
				continue
			}
			return nil, err
		}
		if info.IsDir() {
			return nil, fmt.Errorf("%s is a directory", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		plan.Info = info
		plan.Content = data
		plan.HasManagedBlock = hasManagedAgentsBlock(string(data))
		plans = append(plans, plan)
	}
	return plans, nil
}

func (p agentInstructionFilePlan) DryRunMessage() string {
	if p.Missing {
		return fmt.Sprintf("Would skip %s (file not found)", p.Path)
	}
	if p.HasManagedBlock {
		return fmt.Sprintf("Would update %s (managed block found)", p.Path)
	}
	return fmt.Sprintf("Would update %s (no managed block yet, would append)", p.Path)
}

func mergeAgentsBlock(content string, snippet string) string {
	block := managedAgentsBlock(snippet)
	for _, markers := range agentBlockMarkers() {
		start := strings.Index(content, markers[0])
		if start >= 0 {
			end := strings.Index(content[start:], markers[1])
			if end >= 0 {
				end += start + len(markers[1])
				if strings.HasPrefix(content[end:], "\r\n") {
					end += len("\r\n")
				} else if strings.HasPrefix(content[end:], "\n") {
					end++
				}
				return content[:start] + block + content[end:]
			}
		}
	}

	if strings.TrimSpace(content) == "" {
		return block
	}
	if strings.HasSuffix(content, "\n\n") {
		return content + block
	}
	if strings.HasSuffix(content, "\n") {
		return content + "\n" + block
	}
	return content + "\n\n" + block
}

func hasManagedAgentsBlock(content string) bool {
	for _, markers := range agentBlockMarkers() {
		start := strings.Index(content, markers[0])
		if start >= 0 && strings.Contains(content[start:], markers[1]) {
			return true
		}
	}
	return false
}

func agentBlockMarkers() [][2]string {
	return [][2]string{
		{agentsBeginMarker, agentsEndMarker},
		{legacyAgentsBeginMarker, legacyAgentsEndMarker},
	}
}

func managedAgentsBlock(snippet string) string {
	return agentsBeginMarker + "\n" + strings.TrimRight(snippet, "\n") + "\n" + agentsEndMarker + "\n"
}

func selectedAgentsSnippet(compact bool) string {
	if compact {
		return agentsSnippetCompact
	}
	return agentsSnippet
}
