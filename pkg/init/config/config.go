package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Available config variables
const (
	AppPort = "app.port"
	AppURL  = "app.url"
)

// Config contains and provides the configuration that is required at runtime
type Config interface {
	GetString(string) string
	GetInt(string) int
	GetInt64(string) int64
	GetBool(string) bool
}

// getConfig returns the configuration
func getConfig(path string) (Config, error) {

	// defining that we want to read config from the file named "app" in the provided directory
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
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

func CreateConfig() Config {
	// getting the path of the main file
	ex, err := os.Executable()
	if err != nil {
		log.Panicln(err)
	}

	dir, err := filepath.Abs(filepath.Dir(ex))
	if err != nil {
		log.Panicln(err)
	}

	// loading the config and checking the current directory for an app.yaml file
	cfg, err := getConfig(dir)
	if err != nil {
		log.Panicln(err)
	}

	return cfg
}
