package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/MultivendorEcom/handler"
	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage/postgres"
	"github.com/gorilla/schema"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go:embed assets
var assets embed.FS

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	logger.Info("starting the application")
	env := config.GetString("runtime.environment")
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal("error in reading zap logger")
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	asst, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Println("error while get assets")
	}
	store, err := postgres.NewDBStringFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	st, err := postgres.NewStorage(store)
	if err != nil {
		log.Fatal(err)
	}
	r, err := handler.New(env, config, logger, decoder, asst, st)
	if err != nil {
		log.Fatalf("error in connecting Server : %v", err)
	}
	serverPort := config.GetString("server.port")
	if err := http.ListenAndServe(serverPort, r); err != nil {
		log.Fatalf("error in listen server. %v", err)
	}
}
