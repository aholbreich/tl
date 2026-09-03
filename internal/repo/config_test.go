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
	if cfg.Format != "" {
		t.Errorf("Format: got %q, want empty for a legacy config", cfg.Format)
	}
}

func TestLoadConfigCustomTTL(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(filepath.Join(ledger, TasksDir), 0755); err != nil {
		t.Fatal(err)
	}
	cfgYAML := "format: tl\nversion: 1\ndefault_claim_ttl: 120m\n"
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
	if cfg.Format != LedgerFormat {
		t.Errorf("Format: got %q, want %q", cfg.Format, LedgerFormat)
	}
}

func TestAddFormatMarkerPreservesLegacyConfig(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(ledger, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "# project setting\nversion: 1\ndefault_claim_ttl: 120m\n"
	if err := os.WriteFile(filepath.Join(ledger, ConfigFile), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := AddFormatMarker(ledger); err != nil {
		t.Fatalf("AddFormatMarker: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(ledger, ConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != configIdentity+legacy {
		t.Fatalf("config content:\n%s\nwant:\n%s", got, configIdentity+legacy)
	}

	// Applying the repair again is a no-op.
	if err := AddFormatMarker(ledger); err != nil {
		t.Fatalf("second AddFormatMarker: %v", err)
	}
	again, err := os.ReadFile(filepath.Join(ledger, ConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != string(data) {
		t.Fatalf("second AddFormatMarker changed config:\n%s", again)
	}
}

func TestAddFormatMarkerPreservesYAMLDocumentMarker(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(ledger, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "---\nversion: 1\ndefault_claim_ttl: 120m\n"
	path := filepath.Join(ledger, ConfigFile)
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := AddFormatMarker(ledger); err != nil {
		t.Fatalf("AddFormatMarker: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "---\n" + configIdentity + "version: 1\ndefault_claim_ttl: 120m\n"
	if string(data) != want {
		t.Fatalf("config content:\n%s\nwant:\n%s", data, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("config mode: got %o, want 600", got)
	}
}

func TestAddFormatMarkerRefusesAnotherFormat(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, LedgerDir)
	if err := os.MkdirAll(ledger, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(ledger, ConfigFile)
	original := "format: other\nversion: 1\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := AddFormatMarker(ledger); err == nil {
		t.Fatal("AddFormatMarker should reject a conflicting format")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != original {
		t.Fatalf("conflicting config was modified:\n%s", data)
	}
}
