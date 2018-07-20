package main

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
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

  // Git clone build target.
  log.Infof("Cloning %s...", cfg.BuildTargetUrl)

  // Validate build target's config file.
  log.Infoln("Validating target config file...")

  // Git clone buildpack.
  log.Infof("Cloning %s buildpack...", cfg.Buildpack)

  // Attach buildpack to target.
  log.Infoln("Attaching buildpack to target...")

  // Build augmented target into Docker image.
  log.Infoln("Building target image...")

  // Push Docker image to external repository.
  log.Infoln("Registering target image...")
}