package main

import (
  "fmt"
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/docker"
  "github.com/sweettea-io/build-server/internal/pkg/gogit"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/util/fileutil"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  r "github.com/gomodule/redigo/redis"
)

var cfg *config.Config
var redisPool *r.Pool
var log *logger.Lgr
var dockerClient *docker.Client

func main() {
  // Setup global vars.
  createConfig()
  createRedisPool()
  createLogger()

  // Clone and validate build target repo.
  cloneBuildTarget()
  validateBuildTargetConfig()

  // Clone and attach buildpack to build target.
  cloneBuildpack()
  attachBuildpack()

  // Build and push Docker image.
  createDockerClient()
  buildImage()
  pushImage()

  log.Infof("%s successfully built and pushed.\n", cfg.ImageTag())
}

func createConfig() {
  cfg = config.New()
}

func createRedisPool() {
  redisPool = redis.NewPool(
    cfg.RedisAddress,
    cfg.RedisPassword,
    cfg.RedisPoolMaxActive,
    cfg.RedisPoolMaxIdle,
    cfg.RedisPoolWait,
  )
}

func createLogger() {
  log = &logger.Lgr{
    Logger: logrus.New(),
    Redis: redisPool,
    Stream: cfg.LogStreamKey(),
  }
}

func cloneBuildTarget() {
  // Ensure destination path doesn't already exist.
  err := fileutil.RemoveIfExists(cfg.BuildTargetLocalPath)
  checkErr(err, "Error removing directory")

  log.Infof("Cloning %s...\n", cfg.BuildTargetUrl)

  cloneErr := gogit.CloneAtSha(
    cfg.BuildTargetUrl,
    cfg.BuildTargetSha,
    cfg.BuildTargetLocalPath,
    log.Logger.Out,
  )

  checkErr(cloneErr, "Error cloning target repository")
}

func validateBuildTargetConfig() {
  log.Infoln("Validating target config file...")
  err := targetutil.ValidateConfig(cfg.BuildTargetLocalPath)
  checkErr(err, "Error validating build target's config file")
}

func cloneBuildpack() {
  // Ensure destination path doesn't already exist.
  err := fileutil.RemoveIfExists(cfg.BuildpackLocalPath)
  checkErr(err, "Error removing directory")


  log.Infof("Cloning %s buildpack...\n", cfg.Buildpack)

  cloneErr := gogit.CloneAtSha(
    cfg.BuildpackUrl,
    cfg.BuildpackSha,
    cfg.BuildpackLocalPath,
    log.Logger.Out,
  )

  checkErr(cloneErr, "Error cloning buildpack repository")
}

func attachBuildpack() {
  log.Infoln("Attaching buildpack to target...")

  err := targetutil.AttachBuildpack(
    cfg.Buildpack,
    cfg.BuildpackLocalPath,
    cfg.BuildTargetLocalPath,
    cfg.BuildTargetUid,
  )

  checkErr(err, "Error attaching buildpack to target")
}

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

func buildImage() {
  log.Infoln("Building target image...")

  err := dockerClient.Build(
    cfg.BuildTargetLocalPath,
    cfg.ImageTag(),
    map[string]*string{"TARGET_UID": &cfg.BuildTargetUid}, // build args
  )

  checkErr(err, "Error building Docker image")
}

func pushImage() {
  log.Infoln("Registering target image...")
  err := dockerClient.Push(cfg.ImageTag())
  checkErr(err, "Error registering Docker image")
}

func checkErr(err error, msg string) {
  if err != nil {
    msg = fmt.Sprintf("%s: %s", msg, err.Error())
    log.Errorf(msg)
    panic(msg)
  }
}