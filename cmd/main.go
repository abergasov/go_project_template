package main

import (
	"flag"
	"fmt"
	"go_project_template/internal/config"
	"go_project_template/internal/logger"
	samplerRepo "go_project_template/internal/repository/sampler"
	"go_project_template/internal/routes"
	samplerService "go_project_template/internal/service/sampler"
	"go_project_template/internal/storage/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	flag.Parse()
	appLog, err := logger.NewAppLogger(appHash)
	if err != nil {
		log.Fatalf("unable to create logger: %s", err)
	}
	appLog.Info("app starting", zap.String("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, zap.String("config", *confFile))
	}

	appLog.Info("create storage connections")
	dbConn, err := getDBConnect(&appConf.ConfigDB)
	if err != nil {
		appLog.Fatal("unable to connect to db", err, zap.String("host", appConf.ConfigDB.Address))
	}
	defer func() {
		if err = dbConn.Client().Close(); err != nil {
			appLog.Fatal("unable to close db connection", err)
		}
	}()

	appLog.Info("init repositories")
	repo := samplerRepo.InitRepo(dbConn)

	appLog.Info("init services")
	service := samplerService.InitService(appLog, repo)

	appLog.Info("init http service")
	appHTTPServer := routes.InitAppRouter(appLog, service, fmt.Sprintf(":%d", appConf.AppPort))
	defer func() {
		if err = appHTTPServer.Stop(); err != nil {
			appLog.Fatal("unable to stop http service", err)
		}
	}()
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("unable to start http service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
}

func getDBConnect(cnf *config.DBConf) (*database.DBConnect, error) {
	for i := 0; i < 5; i++ {
		dbConnect, err := database.InitDBConnect(cnf)
		if err == nil {
			return dbConnect, nil
		}
		time.Sleep(time.Duration(i) * time.Second * 5)
	}
	return nil, fmt.Errorf("can't connect to db")
}
