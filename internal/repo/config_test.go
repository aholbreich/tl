package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(filepath.Join(ledger, TasksDir), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ledger, ConfigFile), []byte("version: 1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(ledger)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.DefaultClaimTTL != "60m" {
		t.Errorf("DefaultClaimTTL: got %q, want 60m", cfg.DefaultClaimTTL)
	}
}

func TestLoadConfigCustomTTL(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(filepath.Join(ledger, TasksDir), 0755); err != nil {
		t.Fatal(err)
	}
	cfgYAML := "version: 1\ndefault_claim_ttl: 120m\n"
	if err := os.WriteFile(filepath.Join(ledger, ConfigFile), []byte(cfgYAML), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(ledger)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.DefaultClaimTTL != "120m" {
		t.Errorf("DefaultClaimTTL: got %q, want 120m", cfg.DefaultClaimTTL)
	}
}
