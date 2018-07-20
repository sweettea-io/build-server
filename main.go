package main

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/docker"
  "github.com/sweettea-io/build-server/internal/pkg/gogit"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
)

func main() {
  // Get new config instance.
  cfg := config.New()

  // Create new redis pool.
  pool := redis.NewPool(
    cfg.RedisAddress,
    cfg.RedisPassword,
    cfg.RedisPoolMaxActive,
    cfg.RedisPoolMaxIdle,
    cfg.RedisPoolWait,
  )

  // Create new logger.
  log := logger.New(pool, cfg.LogStreamKey())

  // ------------- CLONE TARGET REPO -------------

  log.Infof("Cloning %s...", cfg.BuildTargetUrl)

  if err := gogit.CloneAtSha(
    cfg.BuildTargetUrl,
    cfg.BuildTargetSha,
    cfg.BuildTargetLocalPath,
    log.Logger.Out,
  ); err != nil {
    log.Errorf("Error cloning target repository: %s", err.Error())
    return
  }

  // ------------- VALIDATE TARGET CONFIG -------------

  log.Infoln("Validating target config file...")

  if err := targetutil.ValidateConfig(cfg.BuildTargetLocalPath); err != nil {
    log.Errorf("Error validating build target's config file: %s", err.Error())
  }

  // ------------- CLONE BUILDPACK REPO -------------

  log.Infof("Cloning %s buildpack...", cfg.Buildpack)

  if err := gogit.CloneAtSha(
    cfg.BuildpackUrl,
    cfg.BuildpackSha,
    cfg.BuildpackLocalPath,
    log.Logger.Out,
  ); err != nil {
    log.Errorf("Error cloning buildpack repository: %s", err.Error())
    return
  }

  // ------------- ATTACH BUILDPACK TO TARGET -------------

  log.Infoln("Attaching buildpack to target...")

  if err := targetutil.AttachBuildpack(
    cfg.Buildpack,
    cfg.BuildpackLocalPath,
    cfg.BuildTargetLocalPath,
    cfg.BuildTargetUid,
  ); err != nil {
    log.Errorf("Error attaching buildpack to target: %s", err.Error())
    return
  }

  // ------------- CREATE DOCKER CLIENT -------------

  if err := docker.Init(
    cfg.DockerHost,
    cfg.DockerAPIVersion,
    map[string]string{},
  ); err != nil {
    log.Errorf("Error initializing new Docker client: %s", err.Error())
    return
  }

  // ------------- BUILD & TAG DOCKER IMAGE -------------

  log.Infoln("Building target image...")

  imageTag := cfg.ImageTag()
  if err := docker.Build(cfg.BuildTargetLocalPath, imageTag); err != nil {
    log.Errorf("Error building Docker image: %s", err.Error())
    return
  }

  // ------------- PUSH DOCKER IMAGE -------------

  log.Infoln("Registering target image...")

  if err := docker.Push(imageTag); err != nil {
    log.Errorf("Error registering Docker image: %s", err.Error())
  }
}