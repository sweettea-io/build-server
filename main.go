package main

import (
  "fmt"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/docker"
  "github.com/sweettea-io/build-server/internal/pkg/gogit"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/targetconfig"
  "github.com/sweettea-io/build-server/internal/pkg/util/buildpack"
  "github.com/sweettea-io/build-server/internal/pkg/util/fileutil"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  r "github.com/gomodule/redigo/redis"
)

// App config
var cfg *config.Config

// Represents build target's config file
var targetConfig *targetconfig.Config

// Buildpack to use in this job.
var bp *buildpack.Buildpack

// Redis pool used for log streaming
var redisPool *r.Pool

// App logger
var log *logger.Lgr

// Docker client with custom build/push functionality
var dockerClient *docker.Client

// main entry point to this job, whose tasks are to:
// (1) Attach a buildpack to a target repo
// (2) Build a Docker image from the result
// (3) Push that Docker image to a registry
func main() {
  // Setup global vars.
  createConfig()
  createRedisPool()
  createLogger()

  // Clone and validate build target repo.
  cloneBuildTarget()
  createBuildTargetConfig()
  validateBuildTargetConfig()

  // Clone and attach buildpack to build target.
  createBuildpack()
  cloneBuildpack()
  attachBuildpack()

  // Build and push Docker image.
  createDockerClient()
  buildImage()
  pushImage()

  log.Infof("%s successfully built and pushed.\n", cfg.ImageTag())
}

// createConfig creates and assigns an app config
// instance to the global `cfg` var.
func createConfig() {
  cfg = config.New()
}

// createRedisPool creates and assigns a new Redis
// pool instance to the global `redisPool` var.
func createRedisPool() {
  redisPool = redis.NewPool(
    cfg.RedisAddress,
    cfg.RedisPassword,
    cfg.RedisPoolMaxActive,
    cfg.RedisPoolMaxIdle,
    cfg.RedisPoolWait,
  )
}

// createLogger creates and assigns a new logger
// instance to the global `log` var.
func createLogger() {
  log = &logger.Lgr{
    Logger: logrus.New(),
    Redis: redisPool,
    Stream: cfg.LogStreamKey,
  }
}

// cloneBuildTarget git-clones the build target repository
// specified in the global app config.
func cloneBuildTarget() {
  // Ensure destination path doesn't already exist.
  err := fileutil.RemoveIfExists(cfg.BuildTargetLocalPath)
  checkErr(err, "Error removing directory")

  log.Infof("Cloning %s...\n", cfg.BuildTargetUrl)

  // Clone build target repo at specified sha.
  cloneErr := gogit.CloneAtSha(
    cfg.BuildTargetUrl,
    cfg.BuildTargetSha,
    cfg.BuildTargetAccessToken,
    cfg.BuildTargetLocalPath,
    log.Logger.Out,
  )

  checkErr(cloneErr, "Error cloning target repository")
}

func createBuildTargetConfig() {
  log.Infoln("Parsing target config file...")

  var err error
  targetConfig, err = targetconfig.New(cfg.BuildTargetLocalPath)

  checkErr(err, "Error parsing build target's config file")
}

// validateBuildTargetConfig ensures the SweetTea config file exists inside the
// cloned build target repository and that it matches the desired key:value structure.
func validateBuildTargetConfig() {
  log.Infoln("Validating target config file...")
  err := targetConfig.Validate()
  checkErr(err, "Error validating build target's config file")
}

func createBuildpack() {
  var err error
  bp, err = cfg.Buildpack(targetConfig)
  checkErr(err, "Error configuring buildpack")
}

// cloneBuildpack git repository.
func cloneBuildpack() {
  // Ensure destination path doesn't already exist.
  err := fileutil.RemoveIfExists(cfg.BuildpackLocalPath)
  checkErr(err, "Error removing directory")

  log.Infof("Cloning %s buildpack...\n", bp.Name)

  // Clone buildpack repo at specified sha.
  cloneErr := gogit.CloneAtSha(
    bp.Url,
    bp.Sha,
    bp.AccessToken,
    cfg.BuildpackLocalPath,
    log.Logger.Out,
  )

  checkErr(cloneErr, "Error cloning buildpack repository")
}

// attachBuildpack moves files/directories from the cloned
// buildpack repository into the cloned build target repository.
func attachBuildpack() {
  log.Infoln("Attaching buildpack to target...")

  err := targetutil.AttachBuildpack(
    bp,
    cfg.BuildpackLocalPath,
    cfg.BuildTargetLocalPath,
    cfg.BuildTargetUid,
  )

  checkErr(err, "Error attaching buildpack to target")
}

// createDockerClient creates and assigns a new Docker
// client instance to the global `dockerClient` var.
func createDockerClient() {
  var err error

  dockerClient, err = docker.New(
    cfg.DockerHost,
    cfg.DockerAPIVersion,
    map[string]string{}, // default headers
    cfg.DockerRegistryUsername,
    cfg.DockerRegistryPassword,
  )

  checkErr(err, "Error initializing new Docker client")
}

// buildImage builds a Docker image from the build target
// after the buildpack has been attached.
func buildImage() {
  log.Infoln("Building target image...")

  err := dockerClient.Build(
    cfg.BuildTargetLocalPath,
    cfg.ImageTag(),
    map[string]*string{"TARGET_UID": &cfg.BuildTargetUid}, // build args
  )

  checkErr(err, "Error building Docker image")
}

// pushImage pushes the Docker image built during
// `buildImage()` to a remote registry.
func pushImage() {
  log.Infoln("Registering target image...")
  err := dockerClient.Push(cfg.ImageTag())
  checkErr(err, "Error registering Docker image")
}

// checkErr checks if an error exists, and if so, does the following:
// (1) constructs the final error message from the provided `msg` and `err.Error()`
// (2) logs the error message to both stderr and the Redis log stream
// (3) panics
func checkErr(err error, msg string) {
  if err != nil {
    msg = fmt.Sprintf("%s: %s", msg, err.Error())
    log.Errorf(msg)
    panic(msg)
  }
}