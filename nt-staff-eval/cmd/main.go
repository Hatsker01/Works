package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Hatsker01/nt-staff-eval/api"
	"github.com/Hatsker01/nt-staff-eval/config"
	"github.com/Hatsker01/nt-staff-eval/pkg/db"
	"github.com/Hatsker01/nt-staff-eval/pkg/logger"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "NT_STAFF")

	defer func(l logger.Logger) {
		err = logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
	)

	conDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	defer conDB.Close()

	server := api.Routers(api.Option{
		Db:     conDB,
		Conf:   cfg,
		Logger: log,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
