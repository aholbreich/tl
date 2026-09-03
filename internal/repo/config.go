package repo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds values read from .tl/config.yaml.
type Config struct {
	Format          string       `yaml:"format"`
	DefaultClaimTTL string       `yaml:"default_claim_ttl"`
	DefaultActor    string       `yaml:"default_actor"`
	Actors          ActorsConfig `yaml:"actors"`
}

// ActorsConfig holds actor-related settings.
type ActorsConfig struct {
	RequireActor bool `yaml:"require_actor"`
}

// LoadConfig reads the config file under ledger. Missing or unparseable values
// fall back to safe defaults (60m TTL).
func LoadConfig(ledger string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(ledger, ConfigFile))
	if err != nil {
		return nil, err
	}
	cfg := &Config{DefaultClaimTTL: "60m"}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.DefaultClaimTTL == "" {
		cfg.DefaultClaimTTL = "60m"
	}
	return cfg, nil
}

// AddFormatMarker adds tl's machine-readable identity to an older config while
// preserving its existing bytes. It is idempotent and refuses to replace a
// marker belonging to another format.
func AddFormatMarker(ledger string) error {
	path := filepath.Join(ledger, ConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var identity struct {
		Format string `yaml:"format"`
	}
	if err := yaml.Unmarshal(data, &identity); err != nil {
		return err
	}
	if identity.Format != "" {
		if identity.Format == LedgerFormat {
			return nil
		}
		return fmt.Errorf("config format is %q, expected %q", identity.Format, LedgerFormat)
	}

	var updated []byte
	if documentMarker := []byte("---\n"); bytes.HasPrefix(data, documentMarker) {
		updated = append(updated, documentMarker...)
		updated = append(updated, configIdentity...)
		updated = append(updated, data[len(documentMarker):]...)
	} else {
		updated = append(updated, configIdentity...)
		updated = append(updated, data...)
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(info.Mode().Perm()); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(updated); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
