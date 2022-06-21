package configs

import (
	"github.com/spf13/viper"
	"os"
)

func Init(path string) error {
	viper.SetConfigFile(path)

	// todo: определить дефолтные настройки
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

type Config struct {
	Addr                                 string
	Endpoint, AccessKey, SecretAccessKey string
	OnecHost, OnecToken                  string
	Production                           bool
	PG                                   string
}

func ReadConfig() Config {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	onecHost := os.Getenv("ONEC_HOST")
	onceToken := os.Getenv("ONEC_TOKEN")
	addr := os.Getenv("ADDR")
	pga := os.Getenv("PG_ADDR")

	production := true
	if _, ok := os.LookupEnv("DEVELOPMENT"); ok {
		production = false
	}

	return Config{
		Addr:            addr,
		Endpoint:        endpoint,
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		OnecHost:        onecHost,
		OnecToken:       onceToken,
		Production:      production,
		PG:              pga,
	}
}
