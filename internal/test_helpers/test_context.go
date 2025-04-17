package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"go_project_template/internal/config"
	"go_project_template/internal/logger"
	"go_project_template/internal/repository"
	samplerService "go_project_template/internal/service/sampler"
	"go_project_template/internal/storage/database"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	Ctx    context.Context
	Cfg    *config.AppConfig
	Logger logger.AppLogger

	Repo *repository.Repo

	ServiceSampler *samplerService.Service
}

func GetClean(t *testing.T) *TestContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	conf := getTestConfig()
	prepareTestDB(ctx, t, &conf.ConfigDB)

	dbConnect, err := database.InitDBConnect(ctx, &conf.ConfigDB, guessMigrationDir(t))
	require.NoError(t, err)
	cleanupDB(t, dbConnect)
	t.Cleanup(func() {
		cancel()
		require.NoError(t, dbConnect.Client().Close())
	})

	appLog := logger.NewAppSLogger()
	// repo init
	repo := repository.InitRepo(dbConnect)

	// service init
	serviceSampler := samplerService.InitService(ctx, appLog, repo)
	return &TestContainer{
		Ctx:            ctx,
		Cfg:            conf,
		Logger:         appLog,
		Repo:           repo,
		ServiceSampler: serviceSampler,
	}
}

func prepareTestDB(ctx context.Context, t *testing.T, cnf *config.DBConf) {
	dbConnect, err := database.InitDBConnect(ctx, &config.DBConf{
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
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
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
	for _, table := range repository.AllTables {
		_, err := connector.Client().Exec(fmt.Sprintf("TRUNCATE %s CASCADE", table))
		require.NoError(t, err)
	}
}
