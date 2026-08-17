package config_test

import (
	"flag"
	"go_project_template/internal/utils/config"
	"os"
	"path/filepath"
	"testing"
)

// FuzzConfigLoad tests YAML config loading with arbitrary input
func FuzzConfigLoad(f *testing.F) {
	f.Add([]byte(`port: 8080`))
	f.Add([]byte(`enabled: true`))
	f.Add([]byte(`values: [1, 2, 3]`))
	f.Add([]byte(`grpc_port: "50051"`))
	f.Add([]byte(``))
	f.Add([]byte(`:`))
	f.Add([]byte(`{invalid`))
	f.Add([]byte(`key: [unclosed`))

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(tmpFile, data, 0600); err != nil {
			return
		}

		type TestConfig struct {
			Port     int      `yaml:"port" default:"8080"`
			Enabled  bool     `yaml:"enabled" default:"false"`
			Values   []int    `yaml:"values"`
			GRPCPort string   `yaml:"grpc_port" default:"50051"`
			LogLevel string   `yaml:"log_level" default:"info"`
			Hosts    []string `yaml:"hosts"`
		}

		// Use a custom FlagSet to avoid parsing os.Args
		fs := flag.NewFlagSet("fuzz", flag.ContinueOnError)
		_ = fs.Parse([]string{}) // Parse empty args

		var cfg TestConfig
		_ = config.Load(&cfg, config.WithYAML(tmpFile), config.WithFlagSet(fs))
	})
}
