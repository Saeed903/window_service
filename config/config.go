package config

import (
	"fmt"
	"os"

	"github.com/saeed903/windows_service/pkg/constants"
	"github.com/spf13/viper"
)

var ConfigPath string

type Config struct {
	FolderWatchPath string `mapstructure:"folderWatchPath"`
}

func InitConfig() (*Config, error) {
	if ConfigPath == "" {
		ConfigPathFromEnv := os.Getenv(constants.ConfigPath)
		if ConfigPathFromEnv != "" {
			ConfigPath = ConfigPathFromEnv
		} else {
			ConfigPath = fmt.Sprint("./config/config.yaml")
		}
 	}

	cfg := &Config{}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(ConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}