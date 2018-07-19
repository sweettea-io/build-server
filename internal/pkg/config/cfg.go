package config

import (
  "fmt"
  "github.com/sweettea-io/envdecode"
)

// Config represents app config populated from environment variables.
type Config struct {
  Debug              bool   `env:"DEBUG,required"`
  Env                string `env:"ENV,required"`
  RedisPoolMaxActive int    `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle   int    `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait      bool   `env:"REDIS_POOL_WAIT,required"`
  RedisAddress       string `env:"REDIS_ADDRESS,required"`
  RedisPassword      string `env:"REDIS_PASSWORD"`
}

// New creates and returns a new Config struct instance populated from environment variables.
func New() *Config {
  // Unmarshal values into a config struct.
  var cfg Config
  if err := envdecode.Decode(&cfg); err != nil {
    panic(fmt.Errorf("Failed to load app config: %s\n", err.Error()))
  }

  // Validate config values.
  validateConfigs(&cfg)

  return &cfg
}

// Validate Config values even further than what has
// already been done during the Decode process.
func validateConfigs(cfg *Config) {}