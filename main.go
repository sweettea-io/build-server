package main

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/docker"
  "github.com/sweettea-io/build-server/internal/pkg/gogit"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  r "github.com/gomodule/redigo/redis"
  "github.com/Sirupsen/logrus"
  "fmt"
)

var cfg *config.Config
var redisPool *r.Pool
var log *logger.Lgr

func main() {
  createConfig()

  createRedisPool()

  createLogger()

  cloneBuildTarget()

  validateBuildTargetConfig()

  cloneBuildpack()

  attachBuildpack()

  createDockerClient()

  buildImage()

  pushImage()
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
  log.Infof("Cloning %s...\n", cfg.BuildTargetUrl)

  err := gogit.CloneAtSha(
    cfg.BuildTargetUrl,
    cfg.BuildTargetSha,
    cfg.BuildTargetLocalPath,
    log.Logger.Out,
  )

  checkErr(err, "Error cloning target repository")
}

func validateBuildTargetConfig() {
  log.Infoln("Validating target config file...")
  err := targetutil.ValidateConfig(cfg.BuildTargetLocalPath)
  checkErr(err, "Error validating build target's config file")
}

func cloneBuildpack() {
  log.Infof("Cloning %s buildpack...\n", cfg.Buildpack)

  err := gogit.CloneAtSha(
    cfg.BuildpackUrl,
    cfg.BuildpackSha,
    cfg.BuildpackLocalPath,
    log.Logger.Out,
  )

  checkErr(err, "Error cloning buildpack repository")
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
  err := docker.Init(
    cfg.DockerHost,
    cfg.DockerAPIVersion,
    map[string]string{},
  )

  checkErr(err, "Error initializing new Docker client")
}

func buildImage() {
  log.Infoln("Building target image...")

  err := docker.Build(
    cfg.BuildTargetLocalPath,
    cfg.ImageTag(),
    map[string]*string{"TARGET_UID": &cfg.BuildTargetUid},
  )

  checkErr(err, "Error building Docker image")
}

func pushImage() {
  log.Infoln("Registering target image...")
  err := docker.Push(cfg.ImageTag())
  checkErr(err, "Error registering Docker image")
}

func checkErr(err error, msg string) {
  if err != nil {
    msg = fmt.Sprintf("%s: %s", msg, err.Error())
    log.Errorf(msg)
    panic(msg)
  }
}