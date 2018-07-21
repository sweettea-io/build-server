package config

import (
  "fmt"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  "github.com/sweettea-io/envdecode"
)

// Config represents app config populated from environment variables.
type Config struct {
  Buildpack              string `env:"BUILDPACK,required"`
  BuildpackAccessToken   string `env:"BUILDPACK_ACCESS_TOKEN"`
  BuildpackLocalPath     string `env:"BUILDPACK_LOCAL_PATH,default=/tmp/buildpack"`
  BuildpackSha           string `env:"BUILDPACK_SHA,required"`
  BuildpackUrl           string `env:"BUILDPACK_URL,required"`
  BuildTargetAccessToken string `env:"BUILD_TARGET_ACCESS_TOKEN"`
  BuildTargetLocalPath   string `env:"BUILD_TARGET_LOCAL_PATH,default=/tmp/target"`
  BuildTargetSha         string `env:"BUILD_TARGET_SHA,required"`
  BuildTargetUid         string `env:"BUILD_TARGET_UID,required"`
  BuildTargetUrl         string `env:"BUILD_TARGET_URL,required"`
  Debug                  bool   `env:"DEBUG,required"`
  DeployUid              string `env:"DEPLOY_UID,required"`
  DockerAPIVersion       string `env:"DOCKER_API_VERSION,default=v1.30"`
  DockerHost             string `env:"DOCKER_HOST,default=unix:///var/run/docker.sock"`
  DockerRegistryOrg      string `env:"DOCKER_REGISTRY_ORG,required"`
  DockerRegistryUsername string `env:"DOCKER_REGISTRY_USERNAME,required"`
  DockerRegistryPassword string `env:"DOCKER_REGISTRY_PASSWORD,required"`
  Env                    string `env:"ENV,required"`
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

// ImageTag returns the tag to be added to the Docker image being built in this job.
func (cfg *Config) ImageTag() string {
  return fmt.Sprintf("%s/%s-%s:%s", cfg.DockerRegistryOrg, cfg.TargetCluster, cfg.BuildTargetUid, cfg.BuildTargetSha)
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
  if !targetutil.IsValidTargetCluster(cfg.TargetCluster) {
    panic(fmt.Errorf(
      "%s is not a valid target cluster. Check 'internal/pkg/util/targetutil/clusters.go' for a list of valid options.\n",
      cfg.TargetCluster,
    ))
  }
}