package config

import (
	"github.com/spf13/viper"
)

// Read loads the application configuration
func Read() (*Config, error) {
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Set default variables if any
	// viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
