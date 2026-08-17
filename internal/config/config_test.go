package config_test

import (
	"go_project_template/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaultConfig(t *testing.T) {
	// given, when
	cfg, err := config.InitConf(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// then
	require.Equal(t, 8000, cfg.AppPort)
	require.Equal(t, true, cfg.EnableTelemetry)
	require.Equal(t, "localhost", cfg.ConfigDB.Address)
	require.Equal(t, "5432", cfg.ConfigDB.Port)
	require.Equal(t, "user", cfg.ConfigDB.User)
	require.Equal(t, "", cfg.ConfigDB.Pass)
	require.Equal(t, "dbname", cfg.ConfigDB.DBName)
	require.Equal(t, 10, cfg.ConfigDB.MaxConnections)
}
