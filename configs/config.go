package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl string `mapstructure:"POSTGRESQL_URL"`
	HttpPort    string `mapstructure:"HTTP_PORT"`
	GrpcPort    string `mapstructure:"GRPC_PORT"`
	GraphqlPort string `mapstructure:"GRAPHQL_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error occured while initialisation configs: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("error occured while initialisation configs: %w", err)
	}
	return config, err
}
