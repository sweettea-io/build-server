package config

import (
  "fmt"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  "github.com/sweettea-io/envdecode"
  "github.com/sweettea-io/build-server/internal/pkg/util/buildpack"
  "github.com/sweettea-io/build-server/internal/pkg/targetconfig"
  "github.com/sweettea-io/build-server/internal/pkg/util/lang"
)

// Config represents app config populated from environment variables
type Config struct {
  BuildpackLocalPath                     string `env:"BUILDPACK_LOCAL_PATH,default=/tmp/buildpack"`
  BuildTargetAccessToken                 string `env:"BUILD_TARGET_ACCESS_TOKEN"`
  BuildTargetLocalPath                   string `env:"BUILD_TARGET_LOCAL_PATH,default=/tmp/target"`
  BuildTargetSha                         string `env:"BUILD_TARGET_SHA,required"`
  BuildTargetUid                         string `env:"BUILD_TARGET_UID,required"`
  BuildTargetUrl                         string `env:"BUILD_TARGET_URL,required"`
  DockerAPIVersion                       string `env:"DOCKER_API_VERSION,default=v1.30"`
  DockerHost                             string `env:"DOCKER_HOST,default=unix:///var/run/docker.sock"`
  DockerRegistryOrg                      string `env:"DOCKER_REGISTRY_ORG,required"`
  DockerRegistryUsername                 string `env:"DOCKER_REGISTRY_USERNAME,required"`
  DockerRegistryPassword                 string `env:"DOCKER_REGISTRY_PASSWORD,required"`
  LogStreamKey                           string `env:"LOG_STREAM_KEY,required"`
  PythonJsonApiBuildpackAccessToken      string `env:"PYTHON_JSON_API_BUILDPACK_ACCESS_TOKEN"`
  PythonJsonApiBuildpackSha              string `env:"PYTHON_JSON_API_BUILDPACK_SHA,required"`
  PythonJsonApiBuildpackUrl              string `env:"PYTHON_JSON_API_BUILDPACK_URL,required"`
  PythonTrainBuildpackAccessToken        string `env:"PYTHON_TRAIN_BUILDPACK_ACCESS_TOKEN"`
  PythonTrainBuildpackSha                string `env:"PYTHON_TRAIN_BUILDPACK_SHA,required"`
  PythonTrainBuildpackUrl                string `env:"PYTHON_TRAIN_BUILDPACK_URL,required"`
  PythonWebsocketApiBuildpackAccessToken string `env:"PYTHON_WEBSOCKET_API_BUILDPACK_ACCESS_TOKEN"`
  PythonWebsocketApiBuildpackSha         string `env:"PYTHON_WEBSOCKET_API_BUILDPACK_SHA,required"`
  PythonWebsocketApiBuildpackUrl         string `env:"PYTHON_WEBSOCKET_API_BUILDPACK_URL,required"`
  RedisPoolMaxActive                     int    `env:"REDIS_POOL_MAX_ACTIVE,required"`
  RedisPoolMaxIdle                       int    `env:"REDIS_POOL_MAX_IDLE,required"`
  RedisPoolWait                          bool   `env:"REDIS_POOL_WAIT,required"`
  RedisAddress                           string `env:"REDIS_ADDRESS,required"`
  RedisPassword                          string `env:"REDIS_PASSWORD"`
  TargetCluster                          string `env:"TARGET_CLUSTER,required"`
}

// ImageTag returns the tag for the Docker image being built in this job.
func (cfg *Config) ImageTag() string {
  return fmt.Sprintf("%s/%s-%s:%s", cfg.DockerRegistryOrg, cfg.TargetCluster, cfg.BuildTargetUid, cfg.BuildTargetSha)
}

func (cfg *Config) Buildpack(targetCfg *targetconfig.Config) (*buildpack.Buildpack, error) {
  // Get buildpack name based on the target cluster value.
  bpName := ""
  switch cfg.TargetCluster {
  case targetutil.Train:
    bpName = targetCfg.Training.Buildpack
  case targetutil.API:
    bpName = targetCfg.Hosting.Buildpack
  }

  // Create and return appropriate Buildpack struct for buildpack name.
  switch bpName {
  case buildpack.PythonTrainBuildpack:
    return &buildpack.Buildpack{
      Name: buildpack.PythonTrainBuildpack,
      Url: cfg.PythonTrainBuildpackUrl,
      Sha: cfg.PythonTrainBuildpackSha,
      AccessToken: cfg.PythonTrainBuildpackAccessToken,
      FileExt: lang.Python.FileExt,
    }, nil
  case buildpack.PythonJsonApiBuildpack:
    return &buildpack.Buildpack{
      Name: buildpack.PythonJsonApiBuildpack,
      Url: cfg.PythonJsonApiBuildpackUrl,
      Sha: cfg.PythonJsonApiBuildpackSha,
      AccessToken: cfg.PythonJsonApiBuildpackAccessToken,
      FileExt: lang.Python.FileExt,
    }, nil
  case buildpack.PythonWebsocketApiBuildpack:
    return &buildpack.Buildpack{
      Name: buildpack.PythonWebsocketApiBuildpack,
      Url: cfg.PythonWebsocketApiBuildpackUrl,
      Sha: cfg.PythonWebsocketApiBuildpackSha,
      AccessToken: cfg.PythonWebsocketApiBuildpackAccessToken,
      FileExt: lang.Python.FileExt,
    }, nil
  default:
    return nil, fmt.Errorf("can't configure buildpack for name: \"%s\"", bpName)
  }
}

// New creates and returns a new config instance populated from environment variables.
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

// validateConfigs applies further validation in relation to what these
// config values actually are. The checks for whether these values exist or not have
// already been taken care of using tags in the `Config` struct definition.
func validateConfigs(cfg *Config) {
  // Ensure TARGET_CLUSTER value is supported.
  if !targetutil.IsValidTargetCluster(cfg.TargetCluster) {
    panic(fmt.Errorf(
      "%s is not a valid target cluster. Check 'internal/pkg/util/targetutil/clusters.go' for a list of valid options.\n",
      cfg.TargetCluster,
    ))
  }
}