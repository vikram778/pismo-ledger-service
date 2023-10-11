package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"pismo-ledger-service/config"
	"pismo-ledger-service/pkg/log"
	"pismo-ledger-service/pkg/postgres"
	"pismo-ledger-service/pkg/repository"
	"pismo-ledger-service/server"
)

func main() {
	log.SetLogLevel()

	log.Info("Starting server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatal("Loading config:", zap.Error(err))
	}

	psqlDB, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Postgresql init: %s", err))
	}
	defer psqlDB.Close()

	log.Info(fmt.Sprintf("PostgreSQL connected: %#v", psqlDB.Stats()))

	s := repository.NewDBOpsRepository(psqlDB)
	h := server.NewServer(cfg, &s)

	log.Fatal("Starting Service Error", zap.Error(http.ListenAndServe(":3000", h.Run())))
}
