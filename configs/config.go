package configs

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl string `mapstructure:"POSTGRESQL_URL"`
	HttpPort    string `mapstructure:"HTTP_PORT"`
	GrpcPort    string `mapstructure:"GRPC_PORT"`
	GraphqlPort string `mapstructure:"GRAPHQL_PORT"`
}

func (c *Config) validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.PostgresUrl, validation.Required, is.URL),
		validation.Field(&c.HttpPort, validation.Required, is.Port, validation.NotIn()),
		validation.Field(&c.GrpcPort, validation.Required, is.Port, validation.NotIn(c.HttpPort)),
		validation.Field(&c.GraphqlPort, validation.Required, is.Port, validation.NotIn(c.HttpPort, c.GrpcPort)),
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		PostgresUrl: viper.GetString("postgresql_url"),
		HttpPort:    viper.GetString("http_port"),
		GrpcPort:    viper.GetString("grpc_port"),
		GraphqlPort: viper.GetString("graphql_port"),
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("error occured while initialisation configs: %w", err)
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("error occured while initialisation configs: %w", err)
	}

	return config, err
}
