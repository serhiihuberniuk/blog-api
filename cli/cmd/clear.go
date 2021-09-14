package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

var allowedTypeFlagValues = map[string]string{
	"user":    "user:*",
	"post":    "post:*",
	"comment": "comment:*",
}

var typeFlagValue string
var allFlagValue bool
var redisAddressFlagValue string

var cacheClearCmd = &cobra.Command{
	Use:   "cache:clear",
	Short: "Clear redis cache storage",
	Long: `cache:clear is a command to clear your Redis cache storage.
You have to provide object-type or use --all flag to execute the command.	
You can change your redis db address by providing --redis-addr flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddressFlagValue,
			DB:       0,
			Password: "",
		})

		err := client.Ping().Err()
		if err != nil {
			return fmt.Errorf("connection to Redis on %v failed: %w", redisAddressFlagValue, err)
		}

		if allFlagValue {
			client.FlushAll()
			log.Println("all cache is cleared")

			return nil
		}

		if typeFlagValue == "" {
			return errors.New("please provide cache type to be cleared")
		}

		value, ok := allowedTypeFlagValues[typeFlagValue]
		if !ok {
			return errors.New(fmt.Sprintf(`object-type "%s" do not exists`, typeFlagValue))
		}

		keys, err := client.Keys(value).Result()
		if err != nil {
			return fmt.Errorf("error occured while finding keys for specified type: %w", err)
		}

		for _, key := range keys {
			err := client.Del(key).Err()
			if err != nil {
				return fmt.Errorf("error occured while removing values from cache: %w", err)
			}
		}
		log.Printf("%s cache is cleared", typeFlagValue)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cacheClearCmd)

	cacheClearCmd.Flags().StringVar(&typeFlagValue, "type", "", "Type of object to clear")
	cacheClearCmd.Flags().BoolVarP(&allFlagValue, "all", "a", false, "Clear all cache")
	cacheClearCmd.Flags().StringVar(&redisAddressFlagValue, "redis-addr", "localhost:6379", "Address of Redis cache storage")

}
