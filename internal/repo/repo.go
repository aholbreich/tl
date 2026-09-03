// Package repo manages the on-disk .tl layout: locating an existing
// ledger and creating a fresh one.
package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	LedgerDir     = ".tl"
	ConfigFile    = "config.yaml"
	TasksDir      = "tasks"
	EventsJournal = "events.jsonl"
	LedgerReadme  = "README.md"
	LedgerFormat  = "tl"
)

// ErrAlreadyInitialized is returned by Init when a ledger directory already
// exists at the target location.
var ErrAlreadyInitialized = errors.New("ledger already exists")

const configIdentity = `# tl task ledger: https://github.com/aholbreich/tl
format: tl
`

const defaultConfig = configIdentity + `version: 1
default_claim_ttl: 60m
id_prefix: task
actors:
  require_actor: true
`

const defaultLedgerReadme = `# tl task ledger

This directory contains the repository's work ledger, maintained with
[tl](https://github.com/aholbreich/tl).

- ` + "`tasks/*.md`" + ` — human-readable tasks
- ` + "`events.jsonl`" + ` — append-only audit history
- ` + "`config.yaml`" + ` — ledger configuration

You can inspect these files without installing ` + "`tl`" + `. When changing the ledger,
use the ` + "`tl`" + ` CLI so its append-only history remains complete.
`

// Init creates the .tl layout under dir and returns the absolute path
// to the created ledger directory. It refuses to touch an existing ledger.
func Init(dir string) (string, error) {
	ledger := filepath.Join(dir, LedgerDir)

	if _, err := os.Stat(ledger); err == nil {
		return "", fmt.Errorf("%w at %s", ErrAlreadyInitialized, ledger)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	if err := os.MkdirAll(filepath.Join(ledger, TasksDir), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(ledger, ConfigFile), []byte(defaultConfig), 0o644); err != nil {
		return "", err
	}
	if err := CreateLedgerReadme(ledger); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(ledger, EventsJournal), nil, 0o644); err != nil {
		return "", err
	}
	return ledger, nil
}

// CreateLedgerReadme adds the human-facing description of a tl ledger. It
// refuses to overwrite an existing file so doctor repairs cannot clobber
// project-specific documentation.
func CreateLedgerReadme(ledger string) error {
	path := filepath.Join(ledger, LedgerReadme)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(defaultLedgerReadme); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	return f.Close()
}
