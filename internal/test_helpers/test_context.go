package testhelpers

import (
	"fmt"
	"go_project_template/internal/config"
	"go_project_template/internal/logger"
	"go_project_template/internal/repository/sampler"
	samplerService "go_project_template/internal/service/sampler"
	"go_project_template/internal/storage/database"
	"os"

	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestContainer struct {
}

func GetClean(t *testing.T) *TestContainer {
	conf := getTestConfig()
	prepareTestDB(t, &conf.ConfigDB)

	dbConnect, err := database.InitDBConnect(&conf.ConfigDB, guessMigrationDir(t))
	require.NoError(t, err)
	cleanupDB(t, dbConnect)
	t.Cleanup(func() {
		require.NoError(t, dbConnect.Client().Close())
	})

	appLog := logger.NewAppSLogger("test")
	// repo init
	repo := sampler.InitRepo(dbConnect)

	// service init
	samplerService.InitService(appLog, repo)
	return &TestContainer{}
}

func prepareTestDB(t *testing.T, cnf *config.DBConf) {
	dbConnect, err := database.InitDBConnect(&config.DBConf{
		Address:        cnf.Address,
		Port:           cnf.Port,
		User:           cnf.User,
		Pass:           cnf.Pass,
		DBName:         "postgres",
		MaxConnections: cnf.MaxConnections,
	}, "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, dbConnect.Client().Close())
	}()
	if _, err = dbConnect.Client().Exec(fmt.Sprintf("CREATE DATABASE %s", cnf.DBName)); !isDatabaseExists(err) {
		require.NoError(t, err)
	}
}

func getTestConfig() *config.AppConfig {
	return &config.AppConfig{
		AppPort: 0,
		ConfigDB: config.DBConf{
			Address:        "localhost",
			Port:           "5449",
			User:           "aHAjeK",
			Pass:           "AOifjwelmc8dw",
			DBName:         "sybill_test",
			MaxConnections: 10,
		},
	}
}

func isDatabaseExists(err error) bool {
	return checkSQLError(err, "42P04")
}

func checkSQLError(err error, code string) bool {
	if err == nil {
		return false
	}
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return string(pqErr.Code) == code
}

func guessMigrationDir(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	res := strings.Split(dir, "/internal")
	return res[0] + "/migrations"
}

func cleanupDB(t *testing.T, connector database.DBConnector) {
	tables := []string{"sampler"}
	for _, table := range tables {
		_, err := connector.Client().Exec(fmt.Sprintf("TRUNCATE %s CASCADE", table))
		require.NoError(t, err)
	}
}
