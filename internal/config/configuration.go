package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Configuration struct {
	Server  ServerConfiguration
	Mongodb MongodbConfiguration
	Redis   RedisConfiguration
	S3      S3StorageConfiguration
	Python  MicroserviceConfiguration
}

// Setup initialize configuration
var (
	// Params ParamsConfiguration
	Config *Configuration
)

// Params = getConfig.Params
func Setup() {
	var configuration *Configuration
	logger := utility.NewLogger()

	viper.SetConfigName("config-sample")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Params = configuration.Params
	Config = configuration
	logger.Info("configurations loading successfully")
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
