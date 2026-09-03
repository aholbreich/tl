package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCreatesIdentifiedLedger(t *testing.T) {
	ledger, err := Init(t.TempDir())
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	cfg, err := LoadConfig(ledger)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Format != LedgerFormat {
		t.Fatalf("ledger format: got %q, want %q", cfg.Format, LedgerFormat)
	}

	data, err := os.ReadFile(filepath.Join(ledger, LedgerReadme))
	if err != nil {
		t.Fatalf("read ledger README: %v", err)
	}
	readme := string(data)
	if !strings.Contains(readme, "# tl task ledger") ||
		!strings.Contains(readme, "https://github.com/aholbreich/tl") {
		t.Fatalf("ledger README does not identify tl:\n%s", readme)
	}
}

func TestCreateLedgerReadmeDoesNotOverwrite(t *testing.T) {
	ledger := t.TempDir()
	path := filepath.Join(ledger, LedgerReadme)
	const existing = "project-specific guide\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CreateLedgerReadme(ledger); err == nil {
		t.Fatal("CreateLedgerReadme should refuse to overwrite an existing file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != existing {
		t.Fatalf("existing README was modified: %q", data)
	}
}
