package main

import (
	"backend/configs"
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

	cfg := configs.ReadConfig()
	rest := Rest(cfg)

	log.Debug().Msg("listen at " + cfg.Addr)
	log.Fatal().Err(http.ListenAndServe(cfg.Addr, rest)).Msg("server stopped")
}
