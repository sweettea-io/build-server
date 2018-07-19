package config

import (
  "fmt"
  "github.com/sweettea-io/envdecode"
)

// Config represents app config populated from environment variables.
type Config struct {
  BuildpackSha           string `env:"BUILDPACK_SHA,required"`
  BuildTargetUid         string `env:"BUILD_TARGET_UID,required"`
  BuildTargetUrl         string `env:"BUILD_TARGET_URL,required"`
  // TODO: remove BuildTargetSha if you can just include the sha inside the BuildTargerUrl
  //BuildTargetSha         string `env:"BUILD_TARGET_SHA,required"`
  Debug                  bool   `env:"DEBUG,required"`
  DeployUid              string `env:"DEPLOY_UID,required"`
  Env                    string `env:"ENV,required"`
  ImageOwner             string `env:"IMAGE_OWNER,required"`
  ImageOwnerPw           string `env:"IMAGE_OWNER_PW,required"`
  RedisPoolMaxActive     int    `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle       int    `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait          bool   `env:"REDIS_POOL_WAIT,required"`
  RedisAddress           string `env:"REDIS_ADDRESS,required"`
  RedisPassword          string `env:"REDIS_PASSWORD"`
  TargetCluster          string `env:"TARGET_CLUSTER,required"`
}

// LogStreamKey returns the redis key of the redis stream used for logging.
func (cfg *Config) LogStreamKey() string {
  return fmt.Sprintf("%s-deploy:%s", cfg.TargetCluster, cfg.DeployUid)
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