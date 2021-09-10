package configs

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUrl             string `mapstructure:"POSTGRESQL_URL"`
	HttpPort                string `mapstructure:"HTTP_PORT"`
	GrpcPort                string `mapstructure:"GRPC_PORT"`
	GraphqlPort             string `mapstructure:"GRAPHQL_PORT"`
	HealthcheckPort         string `mapstructure:"HEALTHCHECK_PORT"`
	PostgresMigrationsPath  string `mapstructure:"POSTGRES_MIGRATIONS_PATH"`
	PrivateKeyFile          string `mapstructure:"PRIVATE_KEY_FILE"`
	RedisAddress            string `mapstructure:"REDIS_ADDRESS"`
	PostgresDatabaseVersion uint   `mapstructure:"POSTGRES_DATABASE_VERSION"`
}

func (c *Config) validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.PostgresUrl, validation.Required, is.URL),
		validation.Field(&c.HttpPort, validation.Required, is.Port),
		validation.Field(&c.GrpcPort, validation.Required, is.Port, validation.NotIn(c.HttpPort)),
		validation.Field(&c.GraphqlPort, validation.Required, is.Port, validation.NotIn(c.HttpPort, c.GrpcPort)),
		validation.Field(&c.HealthcheckPort, validation.Required, is.Port, validation.NotIn(c.HttpPort, c.GrpcPort, c.GraphqlPort)),
		validation.Field(&c.PostgresMigrationsPath, validation.Required),
		validation.Field(&c.PrivateKeyFile, validation.Required),
		validation.Field(&c.RedisAddress, validation.Required),
		validation.Field(&c.PostgresDatabaseVersion, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("api")

	config := &Config{
		PostgresUrl:             viper.GetString("POSTGRESQL_URL"),
		HttpPort:                viper.GetString("HTTP_PORT"),
		GrpcPort:                viper.GetString("GRPC_PORT"),
		GraphqlPort:             viper.GetString("GRAPHQL_PORT"),
		HealthcheckPort:         viper.GetString("HEALTHCHECK_PORT"),
		PostgresMigrationsPath:  viper.GetString("POSTGRES_MIGRATIONS_PATH"),
		PrivateKeyFile:          viper.GetString("PRIVATE_KEY_FILE"),
		RedisAddress:            viper.GetString("REDIS_ADDRESS"),
		PostgresDatabaseVersion: viper.GetUint("POSTGRES_DATABASE_VERSION"),
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("error occurred while initialisation configs: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error occurred while initialisation configs: %w", err)
	}

	return config, nil
}
