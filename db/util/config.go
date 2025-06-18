package util

import "github.com/spf13/viper"

// Config holds the configuration for the database and server.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig loads the configuration from environment variables or a configuration file.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Set default values and override environment variables if they exist
	viper.AutomaticEnv()

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the configuration into the Config struct
	err = viper.Unmarshal(&config)
	return
}
