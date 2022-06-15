package configs

import (
	"github.com/spf13/viper"
)

func Init(path string) error {
	viper.SetConfigFile(path)

	// todo: определить дефолтные настройки
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
