package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.
		With().Caller().
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// - 1c авторизация
	// s3 put get delete
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	onecHost := os.Getenv("ONEC_HOST")
	onceToken := os.Getenv("ONEC_TOKEN")

	r := Rest(cfg{
		endpoint:        endpoint,
		accessKey:       accessKey,
		secretAccessKey: secretAccessKey,
		onecHost:        onecHost,
		onecToken:       onceToken,
	})

	log.Debug().Msg("listen at :" + os.Getenv("HOST"))
	log.Fatal().Err(http.ListenAndServe(":"+os.Getenv("HOST"), r))
}

func logError(w http.ResponseWriter, err error, statusCode int, errorInfo string) {
	log.Err(err).Msg(errorInfo)
	http.Error(w, err.Error(), statusCode)
}
