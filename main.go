package main

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/gogit"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/docker"
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

  log.Infof("Cloning %s...", cfg.BuildTargetUrl)

  // Git clone build target.
  targetCloneErr := gogit.CloneAtSha(
    cfg.BuildTargetUrl,
    cfg.BuildTargetSha,
    cfg.BuildTargetLocalPath,
    log.Logger.Out,
  )

  if targetCloneErr != nil {
    log.Errorf("Error cloning target repository: %s", targetCloneErr.Error())
    return
  }

  log.Infoln("Validating target config file...")

  // Validate build target's config file.

  log.Infof("Cloning %s buildpack...", cfg.Buildpack)

  // Git clone buildpack.
  bpCloneErr := gogit.CloneAtSha(
    cfg.BuildpackUrl,
    cfg.BuildpackSha,
    cfg.BuildpackLocalPath,
    log.Logger.Out,
  )

  if bpCloneErr != nil {
    log.Errorf("Error cloning buildpack repository: %s", bpCloneErr.Error())
    return
  }

  log.Infoln("Attaching buildpack to target...")

  // Attach buildpack to target.

  // Initialize new docker client.
  dockerInitErr := docker.Init(cfg.DockerHost, cfg.DockerAPIVersion, map[string]string{})

  if dockerInitErr != nil {
    log.Errorf("Error initializing new docker client: %s", dockerInitErr.Error())
  }

  log.Infoln("Building target image...")

  // Build augmented target into Docker image.
  dockerBuildErr := docker.Build(cfg.BuildTargetLocalPath, cfg.ImageTag())

  if dockerBuildErr != nil {
    log.Errorf("Error building docker image: %s", dockerBuildErr.Error())
  }

  // Push Docker image to external repository.
  log.Infoln("Registering target image...")
}