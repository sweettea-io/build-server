package config

import (
  "fmt"
  "github.com/sweettea-io/envdecode"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetcluster"
)

// Config represents app config populated from environment variables.
type Config struct {
  Buildpack              string `env:"BUILDPACK,required"`
  BuildpackUrl           string `env:"BUILDPACK_URL,required"`
  BuildpackSha           string `env:"BUILDPACK_SHA,required"`
  BuildpackLocalPath     string `env:"BUILDPACK_LOCAL_PATH,default=/tmp/buildpack"`
  BuildTargetLocalPath   string `env:"BUILD_TARGET_LOCAL_PATH,default=/tmp/target"`
  BuildTargetUid         string `env:"BUILD_TARGET_UID,required"`
  BuildTargetUrl         string `env:"BUILD_TARGET_URL,required"`
  BuildTargetSha         string `env:"BUILD_TARGET_SHA,required"`
  Debug                  bool   `env:"DEBUG,required"`
  DeployUid              string `env:"DEPLOY_UID,required"`
  DockerHost             string `env:"DOCKER_HOST,default=unix:///var/run/docker.sock"`
  DockerAPIVersion       string `env:"DOCKER_API_VERSION,default=v1.30"`
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

func (cfg *Config) ImageTag() string {
  return "<tag>"
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
func validateConfigs(cfg *Config) {
  // Ensure TARGET_CLUSTER value is supported.
  if !targetcluster.IsValidTargetCluster(cfg.TargetCluster) {
    panic(fmt.Errorf(
      "%s is not a valid target cluster. Check 'internal/pkg/util/targetcluster/clusters.go' for a list of valid options.\n",
      cfg.TargetCluster,
    ))
  }
}