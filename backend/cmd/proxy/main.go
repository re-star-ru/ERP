package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.
		With().Caller().
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMilli})

	// - 1c авторизация
	// s3 put get delete
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	onecHost := os.Getenv("ONEC_HOST")
	onceToken := os.Getenv("ONEC_TOKEN")

	production := true
	if _, ok := os.LookupEnv("DEVELOPMENT"); ok {
		production = false
	}

	rest := Rest(cfg{
		endpoint:        endpoint,
		accessKey:       accessKey,
		secretAccessKey: secretAccessKey,
		onecHost:        onecHost,
		onecToken:       onceToken,
		production:      production,
	})

	addr := os.Getenv("ADDR")
	log.Debug().Msg("listen at " + addr)
	log.Fatal().Err(http.ListenAndServe(addr, rest)).Msg("server stopped")
}
