package main

import (
	"backend/configs"
	"backend/internal/app/apiserver"
	"flag"
	"log"

	"github.com/sirupsen/logrus"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
	})

	var path = flag.String("config", "configs/server_config.yml", "path to server main config")
	flag.Parse()

	if err := configs.Init(*path); err != nil {
		log.Fatal(err)
	}

	apiserver.Start()
}
