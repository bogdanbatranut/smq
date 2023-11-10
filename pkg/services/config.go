package services

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	AppURL  = "app.url"
	AppPort = "app.port"
)

type IConfig interface {
	GetString(string) string
	GetInt(string) int
	GetInt64(string) int64
	GetBool(string) bool
}

type ViperConfig struct {
	config IConfig
}

func (cfg ViperConfig) GetString(s string) string {
	return cfg.config.GetString(s)
}

func (cfg ViperConfig) GetInt(s string) int {
	return cfg.config.GetInt(s)
}
func (cfg ViperConfig) GetInt64(s string) int64 {
	return cfg.config.GetInt64(s)
}
func (cfg ViperConfig) GetBool(s string) bool {
	return cfg.config.GetBool(s)
}

func NewViperConfig() (IConfig, error) {
	cfg, err := createViperConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func createViperConfig() (IConfig, error) {

	path, err := getConfigFileDir()
	if err != nil {
		log.Println("No configuration folder found")
	}

	viper.SetConfigName("app")
	viper.AddConfigPath(*path)
	viper.AddConfigPath(".")

	_ = viper.BindEnv(AppURL, "APP_URL")
	_ = viper.BindEnv(AppPort, "APP_PORT")

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	configFileUsed := viper.ConfigFileUsed()
	if len(configFileUsed) == 0 {
		log.Println("no configuration file found")
	} else {
		log.Println("configuration file used")
	}
	return viper.GetViper(), nil
}

func getConfigFileDir() (*string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(ex))
	if err != nil {
		return nil, err
	}
	return &dir, err
}
